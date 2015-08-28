package handler

import (
	"encoding/json"
	"net/http"
	"unicode"
	"unicode/utf8"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"github.com/Lunchr/luncher-api/storage"
	"github.com/deiwin/facebook"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

type HandlerWithUserAndOffer func(w http.ResponseWriter, r *http.Request, user *model.User, offer *model.Offer) *router.HandlerError

// PostOffers handles POST requests to /offers. It stores the offer in the DB and
// sends it to Facebook to be posted on the page's wall at the requested time.
func PostOffers(offersCollection db.Offers, usersCollection db.Users, restaurantsCollection db.Restaurants,
	sessionManager session.Manager, fbAuth facebook.Authenticator, imageStorage storage.Images,
	regionsCollection db.Regions, groupPostsCollection db.OfferGroupPosts) router.Handler {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User) *router.HandlerError {
		restaurant, err := restaurantsCollection.GetID(user.RestaurantIDs[0])
		if err != nil {
			return router.NewHandlerError(err, "Couldn't find a restaurant related to this user", http.StatusInternalServerError)
		}
		offerPOST, err := parseOffer(r, restaurant)
		if err != nil {
			return router.NewHandlerError(err, "Failed to parse the offer", http.StatusBadRequest)
		}
		offer, err := model.MapOfferPOSTToOffer(offerPOST, getImageDataToChecksumMapper(imageStorage))
		if err != nil {
			return router.NewHandlerError(err, "Failed to map the offer to the internal representation", http.StatusInternalServerError)
		}
		offers, err := offersCollection.Insert(offer)
		if err != nil {
			return router.NewHandlerError(err, "Failed to store the offer in the DB", http.StatusInternalServerError)
		}

		date := model.DateFromTime(offer.FromTime)
		handlerErr := updateGroupPostForDate(date, user, restaurant, offersCollection, regionsCollection, groupPostsCollection, fbAuth)
		if handlerErr != nil {
			return handlerErr
		}

		offerJSON, handlerError := mapOfferToJSON(offers[0], imageStorage)
		if handlerError != nil {
			return handlerError
		}
		return writeJSON(w, offerJSON)
	}
	return checkLogin(sessionManager, usersCollection, handler)
}

// PutOffers handles PUT requests to /offers. It updates the offer in the DB and
// updates the related Facebook post.
func PutOffers(offersCollection db.Offers, usersCollection db.Users, restaurantsCollection db.Restaurants,
	sessionManager session.Manager, fbAuth facebook.Authenticator, imageStorage storage.Images,
	regionsCollection db.Regions, groupPostsCollection db.OfferGroupPosts) router.HandlerWithParams {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User, currentOffer *model.Offer) *router.HandlerError {
		restaurant, err := restaurantsCollection.GetID(user.RestaurantIDs[0])
		if err != nil {
			return router.NewHandlerError(err, "Couldn't find the restaurant this offer belongs to", http.StatusInternalServerError)
		}
		offerPOST, err := parseOffer(r, restaurant)
		if err != nil {
			return router.NewHandlerError(err, "Failed to parse the offer", http.StatusBadRequest)
		}
		// If the image_data field isn't set, the image field of offer also doesn't get set and
		// therefore the update won't affect the stored image. This in turn means, that currently
		// there's no way to update the offer to remove an image.
		offer, err := model.MapOfferPOSTToOffer(offerPOST, getImageDataToChecksumMapper(imageStorage))
		if err != nil {
			return router.NewHandlerError(err, "Failed to map the offer to the internal representation", http.StatusInternalServerError)
		}
		err = offersCollection.UpdateID(currentOffer.ID, offer)
		if err != nil {
			return router.NewHandlerError(err, "Failed to update the offer in DB", http.StatusInternalServerError)
		}
		offer.ID = currentOffer.ID

		date := model.DateFromTime(offer.FromTime)
		handlerErr := updateGroupPostForDate(date, user, restaurant, offersCollection, regionsCollection, groupPostsCollection, fbAuth)
		if handlerErr != nil {
			return handlerErr
		}

		offerJSON, handlerError := mapOfferToJSON(offer, imageStorage)
		if handlerError != nil {
			return handlerError
		}
		return writeJSON(w, offerJSON)
	}

	return checkLoginWithParams(sessionManager, usersCollection, forOffer(offersCollection, handler))
}

// DeleteOffers handles DELETE requests to /offers. It deletes the offer from the DB and
// deletes the related Facebook post.
func DeleteOffers(offersCollection db.Offers, usersCollection db.Users, sessionManager session.Manager,
	fbAuth facebook.Authenticator, restaurantsCollection db.Restaurants, regionsCollection db.Regions,
	groupPostsCollection db.OfferGroupPosts) router.HandlerWithParams {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User, currentOffer *model.Offer) *router.HandlerError {
		restaurant, err := restaurantsCollection.GetID(user.RestaurantIDs[0])
		if err != nil {
			return router.NewHandlerError(err, "Couldn't find the restaurant this offer belongs to", http.StatusInternalServerError)
		}
		if err = offersCollection.RemoveID(currentOffer.ID); err != nil {
			return router.NewHandlerError(err, "Failed to delete the offer from DB", http.StatusInternalServerError)
		}

		date := model.DateFromTime(currentOffer.FromTime)
		handlerErr := updateGroupPostForDate(date, user, restaurant, offersCollection, regionsCollection, groupPostsCollection, fbAuth)
		if handlerErr != nil {
			return handlerErr
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
			return router.NewStringHandlerError("invalid offer ID", "", http.StatusBadRequest)
		}
		id := bson.ObjectIdHex(idString)
		offer, err := offersCollection.GetID(id)
		if err != nil {
			return router.NewHandlerError(err, "Couldn't find an offer with this ID", http.StatusNotFound)
		}
		return handler(w, r, user, offer)
	}
}

// postOfferToFB forms a post, sends it to FB and returns the post's FB ID
func postOfferToFB(offer model.Offer, user *model.User, restaurant *model.Restaurant, api facebook.API) (string, *router.HandlerError) {
	if restaurant.FacebookPageID == "" {
		return "", nil
	}
	message := formFBOfferMessage(&offer)
	post, err := api.PagePublish(user.Session.FacebookPageToken, restaurant.FacebookPageID, message)
	if err != nil {
		return "", router.NewHandlerError(err, "Failed to post the offer to Facebook", http.StatusBadGateway)
	}
	return post.ID, nil
}

func capitalizeString(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

func parseOffer(r *http.Request, restaurant *model.Restaurant) (*model.OfferPOST, error) {
	var offer model.OfferPOST
	err := json.NewDecoder(r.Body).Decode(&offer)
	if err != nil {
		return nil, err
	}
	offer.Restaurant = model.OfferRestaurant{
		ID:       restaurant.ID,
		Name:     restaurant.Name,
		Region:   restaurant.Region,
		Address:  restaurant.Address,
		Location: restaurant.Location,
		Phone:    restaurant.Phone,
	}
	return &offer, nil
}

func getImageDataToChecksumMapper(imageStorage storage.Images) func(string) (string, error) {
	return func(imageData string) (string, error) {
		if imageData == "" {
			return "", nil
		}
		return parseAndStoreImage(imageData, imageStorage)
	}
}

func parseAndStoreImage(imageDataURL string, imageStorage storage.Images) (string, error) {
	imageChecksum, err := imageStorage.ChecksumDataURL(imageDataURL)
	if err != nil {
		return "", err
	}
	alreadyStored, err := imageStorage.HasChecksum(imageChecksum)
	if err != nil {
		return "", err
	} else if alreadyStored {
		return imageChecksum, nil
	}
	if err = imageStorage.StoreDataURL(imageDataURL); err != nil {
		return "", err
	}
	return imageChecksum, nil
}
