package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/facebook"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"github.com/Lunchr/luncher-api/storage"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

type HandlerWithRestaurantAndOffer func(http.ResponseWriter, *http.Request, *model.User, *model.Restaurant, *model.Offer) *router.HandlerError

// PostOffers handles POST requests to /offers. It stores the offer in the DB and
// sends it to Facebook to be posted on the page's wall at the requested time.
func PostOffers(offers db.Offers, users db.Users, restaurants db.Restaurants, sessionManager session.Manager,
	imageStorage storage.Images, facebookPost facebook.Post, regions db.Regions) router.HandlerWithParams {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User, restaurant *model.Restaurant) *router.HandlerError {
		offerPOST, err := parseOffer(r, restaurant)
		if err != nil {
			return router.NewHandlerError(err, "Failed to parse the offer", http.StatusBadRequest)
		}
		offer, err := model.MapOfferPOSTToOffer(offerPOST, getImageDataToChecksumMapper(imageStorage))
		if err != nil {
			return router.NewHandlerError(err, "Failed to map the offer to the internal representation", http.StatusInternalServerError)
		}
		offers, err := offers.Insert(offer)
		if err != nil {
			return router.NewHandlerError(err, "Failed to store the offer in the DB", http.StatusInternalServerError)
		}

		location, handlerErr := getLocationForRestaurant(restaurant, regions)
		if handlerErr != nil {
			return handlerErr
		}
		date := model.DateFromTime(offer.FromTime, location)
		handlerErr = facebookPost.Update(date, user, restaurant)
		if handlerErr != nil {
			return handlerErr
		}

		offerJSON, handlerError := mapOfferToJSON(offers[0], imageStorage)
		if handlerError != nil {
			return handlerError
		}
		return writeJSON(w, offerJSON)
	}
	return forRestaurant(sessionManager, users, restaurants, handler)
}

// PutOffers handles PUT requests to /offers. It updates the offer in the DB and
// updates the related Facebook post.
func PutOffers(offers db.Offers, users db.Users, restaurants db.Restaurants, sessionManager session.Manager,
	imageStorage storage.Images, facebookPost facebook.Post, regions db.Regions) router.HandlerWithParams {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User, restaurant *model.Restaurant, currentOffer *model.Offer) *router.HandlerError {
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
		err = offers.UpdateID(currentOffer.ID, offer)
		if err != nil {
			return router.NewHandlerError(err, "Failed to update the offer in DB", http.StatusInternalServerError)
		}
		offer.ID = currentOffer.ID

		location, handlerErr := getLocationForRestaurant(restaurant, regions)
		if handlerErr != nil {
			return handlerErr
		}
		date := model.DateFromTime(offer.FromTime, location)
		handlerErr = facebookPost.Update(date, user, restaurant)
		if handlerErr != nil {
			return handlerErr
		}
		previousDate := model.DateFromTime(currentOffer.FromTime, location)
		if previousDate != date {
			handlerErr = facebookPost.Update(previousDate, user, restaurant)
			if handlerErr != nil {
				return handlerErr
			}
		}

		offerJSON, handlerError := mapOfferToJSON(offer, imageStorage)
		if handlerError != nil {
			return handlerError
		}
		return writeJSON(w, offerJSON)
	}

	return forRestaurantWithParams(sessionManager, users, restaurants, forOffer(offers, handler))
}

// DeleteOffers handles DELETE requests to /offers. It deletes the offer from the DB and
// deletes the related Facebook post.
func DeleteOffers(offers db.Offers, users db.Users, sessionManager session.Manager, restaurants db.Restaurants,
	facebookPost facebook.Post, regions db.Regions) router.HandlerWithParams {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User, restaurant *model.Restaurant, currentOffer *model.Offer) *router.HandlerError {
		if err := offers.RemoveID(currentOffer.ID); err != nil {
			return router.NewHandlerError(err, "Failed to delete the offer from DB", http.StatusInternalServerError)
		}

		location, handlerErr := getLocationForRestaurant(restaurant, regions)
		if handlerErr != nil {
			return handlerErr
		}
		date := model.DateFromTime(currentOffer.FromTime, location)
		handlerErr = facebookPost.Update(date, user, restaurant)
		if handlerErr != nil {
			return handlerErr
		}

		w.WriteHeader(http.StatusOK)
		return nil
	}
	return forRestaurantWithParams(sessionManager, users, restaurants, forOffer(offers, handler))
}

func forOffer(offersCollection db.Offers, handler HandlerWithRestaurantAndOffer) HandlerWithParamsWithRestaurant {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user *model.User, restaurant *model.Restaurant) *router.HandlerError {
		idString := ps.ByName("id")
		if !bson.IsObjectIdHex(idString) {
			return router.NewStringHandlerError("invalid offer ID", "", http.StatusBadRequest)
		}
		id := bson.ObjectIdHex(idString)
		offer, err := offersCollection.GetID(id)
		if err != nil {
			return router.NewHandlerError(err, "Couldn't find an offer with this ID", http.StatusNotFound)
		}
		return handler(w, r, user, restaurant, offer)
	}
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

func getLocationForRestaurant(restaurant *model.Restaurant, regions db.Regions) (*time.Location, *router.HandlerError) {
	region, err := regions.GetName(restaurant.Region)
	if err != nil {
		return nil, router.NewHandlerError(err, "Failed to find the restaurant's region", http.StatusInternalServerError)
	}
	location, err := time.LoadLocation(region.Location)
	if err != nil {
		return nil, router.NewHandlerError(err, "Failed to load region's location", http.StatusInternalServerError)
	}
	return location, nil
}
