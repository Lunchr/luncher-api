package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/deiwin/facebook"
	"github.com/deiwin/imstor"
	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	. "github.com/deiwin/luncher-api/router"
	"github.com/deiwin/luncher-api/session"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

// Offers handles GET requests to /offers. It returns all current day's offers.
func Offers(offersCollection db.Offers, regionsCollection db.Regions, imageStorage imstor.Storage) Handler {
	return func(w http.ResponseWriter, r *http.Request) *HandlerError {
		regionName := r.FormValue("region")
		if regionName == "" {
			return &HandlerError{errors.New("Region not specified for GET /offers"), "Please specify a region", http.StatusBadRequest}
		}
		region, err := regionsCollection.Get(regionName)
		if err != nil {
			return &HandlerError{err, "", http.StatusInternalServerError}
		}
		loc, err := time.LoadLocation(region.Location)
		if err != nil {
			return &HandlerError{err, "", http.StatusInternalServerError}
		}
		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		endTime := startTime.AddDate(0, 0, 1)
		offers, err := offersCollection.Get(regionName, startTime, endTime)
		if err != nil {
			return &HandlerError{err, "", http.StatusInternalServerError}
		}
		for _, offer := range offers {
			if offer.Image != "" {
				imagePath, err := imageStorage.PathForSize(offer.Image, "large")
				if err != nil {
					return &HandlerError{err, "", http.StatusInternalServerError}
				}
				offer.Image = path.Join("images", imagePath)
			}
		}
		return writeJSON(w, offers)
	}
}

// PostOffers handles POST requests to /offers. It stores the offer in the DB and
// sends it to Facebook to be posted on the page's wall at the requested time.
func PostOffers(offersCollection db.Offers, usersCollection db.Users, restaurantsCollection db.Restaurants,
	sessionManager session.Manager, fbAuth facebook.Authenticator, imageStorage imstor.Storage) Handler {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User) *HandlerError {
		api := fbAuth.APIConnection(&user.Session.FacebookUserToken)
		restaurant, err := restaurantsCollection.GetByID(user.RestaurantID)
		if err != nil {
			return &HandlerError{err, "", http.StatusInternalServerError}
		}
		offer, err := parseOffer(r, restaurant)
		if err != nil {
			return &HandlerError{err, "", http.StatusBadRequest}
		}
		if offer.Image != "" {
			imageChecksum, err := parseAndStoreImage(offer.Image, imageStorage)
			if err != nil {
				return &HandlerError{err, "", http.StatusInternalServerError}
			}
			offer.Image = imageChecksum
		}
		message := formFBOfferMessage(*offer)
		post, err := api.PagePublish(user.Session.FacebookPageToken, user.FacebookPageID, message)
		if err != nil {
			return &HandlerError{err, "", http.StatusBadGateway}
		}
		offer.FBPostID = post.ID
		offers, err := offersCollection.Insert(offer)
		if err != nil {
			return &HandlerError{err, "", http.StatusInternalServerError}
		}

		return writeJSON(w, offers[0])
	}
	return checkLogin(sessionManager, usersCollection, handler)
}

// PutOffers handles PUT requests to /offers. It updates the offer in the DB and
// updates the related Facebook post.
func PutOffers(offersCollection db.Offers, usersCollection db.Users, restaurantsCollection db.Restaurants,
	sessionManager session.Manager, fbAuth facebook.Authenticator, imageStorage imstor.Storage) HandlerWithParams {
	handler := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user *model.User) *HandlerError {
		idString := ps.ByName("id")
		if !bson.IsObjectIdHex(idString) {
			err := errors.New("PUT /offers contained an invalid id")
			return &HandlerError{err, "", http.StatusBadRequest}
		}
		id := bson.ObjectIdHex(idString)
		currentOffer, err := offersCollection.GetByID(id)
		if err != nil {
			return &HandlerError{err, "", http.StatusBadRequest}
		}
		api := fbAuth.APIConnection(&user.Session.FacebookUserToken)
		restaurant, err := restaurantsCollection.GetByID(user.RestaurantID)
		if err != nil {
			return &HandlerError{err, "", http.StatusInternalServerError}
		}
		offer, err := parseOffer(r, restaurant)
		if err != nil {
			return &HandlerError{err, "", http.StatusBadRequest}
		}
		if changed, err := imageChanged(currentOffer.Image, offer.Image, imageStorage); err != nil {
			return &HandlerError{err, "", http.StatusInternalServerError}
		} else if changed {
			imageChecksum, err := parseAndStoreImage(offer.Image, imageStorage)
			if err != nil {
				return &HandlerError{err, "", http.StatusInternalServerError}
			}
			offer.Image = imageChecksum
		}
		if currentOffer.FBPostID != "" {
			err = api.PostDelete(user.Session.FacebookPageToken, currentOffer.FBPostID)
			if err != nil {
				return &HandlerError{err, "", http.StatusBadGateway}
			}
		}
		message := formFBOfferMessage(*offer)
		post, err := api.PagePublish(user.Session.FacebookPageToken, user.FacebookPageID, message)
		if err != nil {
			return &HandlerError{err, "", http.StatusBadGateway}
		}
		offer.FBPostID = post.ID
		err = offersCollection.UpdateID(id, offer)
		if err != nil {
			return &HandlerError{err, "", http.StatusInternalServerError}
		}
		offer.ID = id

		return writeJSON(w, offer)
	}

	return checkLoginWithParams(sessionManager, usersCollection, handler)
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
		Name:   restaurant.Name,
		Region: restaurant.Region,
	}
	return &offer, nil
}

func parseAndStoreImage(imageDataURL string, imageStorage imstor.Storage) (string, error) {
	imageChecksum, err := imageStorage.ChecksumDataURL(imageDataURL)
	if err != nil {
		return "", err
	}
	if err = imageStorage.StoreDataURL(imageDataURL); err != nil {
		return "", err
	}
	return imageChecksum, nil
}

func imageChanged(currentImage, putImage string, imageStorage imstor.Storage) (bool, error) {
	if currentImage != "" {
		currentImagePath, err := imageStorage.PathForSize(currentImage, "large")
		if err != nil {
			return false, err
		}
		if putImage != currentImagePath {
			return true, nil
		}
	} else if putImage != "" {
		return true, nil
	}
	return false, nil
}
