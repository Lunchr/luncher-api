package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/deiwin/facebook"
	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"github.com/deiwin/luncher-api/router"
	"github.com/deiwin/luncher-api/session"
	"github.com/deiwin/luncher-api/storage"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

type HandlerWithUserAndOffer func(w http.ResponseWriter, r *http.Request, user *model.User, offer *model.Offer) *router.HandlerError

// PostOffers handles POST requests to /offers. It stores the offer in the DB and
// sends it to Facebook to be posted on the page's wall at the requested time.
func PostOffers(offersCollection db.Offers, usersCollection db.Users, restaurantsCollection db.Restaurants,
	sessionManager session.Manager, fbAuth facebook.Authenticator, imageStorage storage.Images) router.Handler {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User) *router.HandlerError {
		api := fbAuth.APIConnection(&user.Session.FacebookUserToken)
		restaurant, err := restaurantsCollection.GetID(user.RestaurantID)
		if err != nil {
			return &router.HandlerError{err, "Couldn't find a restaurant related to this user", http.StatusInternalServerError}
		}
		offer, err := parseOffer(r, restaurant)
		if err != nil {
			return &router.HandlerError{err, "Failed to parse the offer", http.StatusBadRequest}
		}
		if offer.Image != "" {
			imageChecksum, err := parseAndStoreImage(offer.Image, imageStorage)
			if err != nil {
				return &router.HandlerError{err, "Failed to store the image", http.StatusInternalServerError}
			}
			offer.Image = imageChecksum
		}
		message := formFBOfferMessage(*offer)
		post, err := api.PagePublish(user.Session.FacebookPageToken, user.FacebookPageID, message)
		if err != nil {
			return &router.HandlerError{err, "Failed to post the offer to Facebook", http.StatusBadGateway}
		}
		offer.FBPostID = post.ID
		offers, err := offersCollection.Insert(offer)
		if err != nil {
			return &router.HandlerError{err, "Failed to store the offer in the DB", http.StatusInternalServerError}
		}

		return writeJSON(w, offers[0])
	}
	return checkLogin(sessionManager, usersCollection, handler)
}

// PutOffers handles PUT requests to /offers. It updates the offer in the DB and
// updates the related Facebook post.
func PutOffers(offersCollection db.Offers, usersCollection db.Users, restaurantsCollection db.Restaurants,
	sessionManager session.Manager, fbAuth facebook.Authenticator, imageStorage storage.Images) router.HandlerWithParams {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User, currentOffer *model.Offer) *router.HandlerError {
		api := fbAuth.APIConnection(&user.Session.FacebookUserToken)
		restaurant, err := restaurantsCollection.GetID(user.RestaurantID)
		if err != nil {
			return &router.HandlerError{err, "Couldn't find the restaurant this offer belongs to", http.StatusInternalServerError}
		}
		offer, err := parseOffer(r, restaurant)
		if err != nil {
			return &router.HandlerError{err, "Failed to parse the offer", http.StatusBadRequest}
		}
		if changed, err := imageChanged(currentOffer.Image, offer.Image, imageStorage); err != nil {
			return err
		} else if changed {
			imageChecksum, err := parseAndStoreImage(offer.Image, imageStorage)
			if err != nil {
				return &router.HandlerError{err, "Failed to store the image", http.StatusInternalServerError}
			}
			offer.Image = imageChecksum
		} else {
			offer.Image = currentOffer.Image
		}
		if currentOffer.FBPostID != "" {
			err = api.PostDelete(user.Session.FacebookPageToken, currentOffer.FBPostID)
			if err != nil {
				return &router.HandlerError{err, "Failed to delete the current post from Facebook", http.StatusBadGateway}
			}
		}
		message := formFBOfferMessage(*offer)
		post, err := api.PagePublish(user.Session.FacebookPageToken, user.FacebookPageID, message)
		if err != nil {
			return &router.HandlerError{err, "Failed to post the offer to Facebook", http.StatusBadGateway}
		}
		offer.FBPostID = post.ID
		err = offersCollection.UpdateID(currentOffer.ID, offer)
		if err != nil {
			return &router.HandlerError{err, "Failed to update the offer in DB", http.StatusInternalServerError}
		}
		offer.ID = currentOffer.ID

		return writeJSON(w, offer)
	}

	return checkLoginWithParams(sessionManager, usersCollection, forOffer(offersCollection, handler))
}

// DeleteOffers handles DELETE requests to /offers. It deletes the offer from the DB and
// deletes the related Facebook post.
func DeleteOffers(offersCollection db.Offers, usersCollection db.Users, sessionManager session.Manager, fbAuth facebook.Authenticator) router.HandlerWithParams {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User, currentOffer *model.Offer) *router.HandlerError {
		fbAPI := fbAuth.APIConnection(&user.Session.FacebookUserToken)
		if currentOffer.FBPostID != "" {
			err := fbAPI.PostDelete(user.Session.FacebookPageToken, currentOffer.FBPostID)
			if err != nil {
				return &router.HandlerError{err, "Failed to delete the current post from Facebook", http.StatusBadGateway}
			}
		}
		err := offersCollection.RemoveID(currentOffer.ID)
		if err != nil {
			return &router.HandlerError{err, "Failed to delete the offer from DB", http.StatusInternalServerError}
		}
		w.WriteHeader(http.StatusOK)
		return nil
	}
	return checkLoginWithParams(sessionManager, usersCollection, forOffer(offersCollection, handler))
}

func forOffer(offersCollection db.Offers, handler HandlerWithUserAndOffer) HandlerWithParamsWithUser {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user *model.User) *router.HandlerError {
		idString := ps.ByName("id")
		if !bson.IsObjectIdHex(idString) {
			err := errors.New("Invalid offer ID")
			return &router.HandlerError{err, "", http.StatusBadRequest}
		}
		id := bson.ObjectIdHex(idString)
		offer, err := offersCollection.GetID(id)
		if err != nil {
			return &router.HandlerError{err, "Couldn't find an offer with this ID", http.StatusNotFound}
		}
		return handler(w, r, user, offer)
	}
}

func formFBOfferMessage(o model.Offer) string {
	ingredients := strings.Join(o.Ingredients, ", ")
	capitalizedIngredients := capitalizeString(ingredients)
	return fmt.Sprintf("%s - %s", o.Title, capitalizedIngredients)
}

func capitalizeString(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

func parseOffer(r *http.Request, restaurant *model.Restaurant) (*model.Offer, error) {
	var offer model.Offer
	err := json.NewDecoder(r.Body).Decode(&offer)
	if err != nil {
		return nil, err
	}
	offer.Restaurant = model.OfferRestaurant{
		Name:     restaurant.Name,
		Region:   restaurant.Region,
		Address:  restaurant.Address,
		Location: restaurant.Location,
	}
	return &offer, nil
}

func parseAndStoreImage(imageDataURL string, imageStorage storage.Images) (string, error) {
	imageChecksum, err := imageStorage.ChecksumDataURL(imageDataURL)
	if err != nil {
		return "", err
	}
	if err = imageStorage.StoreDataURL(imageDataURL); err != nil {
		return "", err
	}
	return imageChecksum, nil
}

func imageChanged(currentImage, putImage string, imageStorage storage.Images) (bool, *router.HandlerError) {
	if currentImage == "" {
		return putImage != "", nil
	}
	currentImagePath, err := imageStorage.PathForLarge(currentImage)
	if err != nil {
		return false, &router.HandlerError{err, "Failed to find the current image for this offer", http.StatusInternalServerError}
	}
	return putImage != currentImagePath, nil
}
