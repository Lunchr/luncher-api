package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"github.com/Lunchr/luncher-api/storage"
)

// Restaurants returns a list of all restaurants
func Restaurants(restaurantsCollection db.Restaurants) router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		restaurants, err := restaurantsCollection.Get()
		if err != nil {
			return router.NewHandlerError(err, "", http.StatusInternalServerError)
		}
		return writeJSON(w, restaurants)
	}
}

// Restaurant returns a router.Handler that returns the restaurant information for the
// restaurant linked to the currently logged in user
func Restaurant(c db.Restaurants, sessionManager session.Manager, users db.Users) router.Handler {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User) *router.HandlerError {
		restaurant, err := c.GetID(user.RestaurantID)
		if err != nil {
			return router.NewHandlerError(err, "Failed to find the restaurant connected to this user", http.StatusInternalServerError)
		}
		return writeJSON(w, restaurant)
	}
	return checkLogin(sessionManager, users, handler)
}

// Restaurant returns a router.Handler that returns the restaurant information for the
// restaurant linked to the currently logged in user
func PostRestaurants(c db.Restaurants, sessionManager session.Manager, users db.Users) router.Handler {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User) *router.HandlerError {
		restaurantRegistration, err := parseRestaurantRegistration(r)
		if err != nil {
			return router.NewHandlerError(err, "Failed to parse the restaurant", http.StatusBadRequest)
		}
		insertedRestaurants, err := c.Insert(&restaurantRegistration.Restaurant)
		if err != nil {
			return router.NewHandlerError(err, "Failed to store the restaurant in the DB", http.StatusInternalServerError)
		}
		var insertedRestaurant = insertedRestaurants[0]
		user.RestaurantID = insertedRestaurant.ID
		if restaurantRegistration.PageID != "" {
			user.FacebookPageID = restaurantRegistration.PageID
		}
		err = users.Update(user.FacebookUserID, user)
		if err != nil {
			// TODO: revert the restaurant insertion we just did?
			return router.NewHandlerError(err, "Failed to store the restaurant in the DB", http.StatusInternalServerError)
		}
		return writeJSON(w, insertedRestaurant)
	}
	return checkLogin(sessionManager, users, handler)
}

// RestaurantOffers returns all upcoming offers for the restaurant linked to the
// currently logged in user
func RestaurantOffers(restaurants db.Restaurants, sessionManager session.Manager, users db.Users, offers db.Offers, imageStorage storage.Images) router.Handler {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User) *router.HandlerError {
		restaurant, err := restaurants.GetID(user.RestaurantID)
		if err != nil {
			return router.NewHandlerError(err, "Failed to find the restaurant connected to this user", http.StatusInternalServerError)
		}
		offers, err := offers.GetForRestaurant(restaurant.Name, time.Now())
		if err != nil {
			return router.NewHandlerError(err, "Failed to find upcoming offers for this restaurant", http.StatusInternalServerError)
		}
		for _, offer := range offers {
			if offer.Image != "" {
				offer.Image, err = imageStorage.PathForLarge(offer.Image)
				if err != nil {
					return router.NewHandlerError(err, "Failed to find an image of an offer", http.StatusInternalServerError)
				}
			}
		}
		return writeJSON(w, offers)
	}
	return checkLogin(sessionManager, users, handler)
}

type RestaurantRegistration struct {
	model.Restaurant
	PageID string `json:"pageID,omitempty"`
}

func parseRestaurantRegistration(r *http.Request) (*RestaurantRegistration, error) {
	var restaurantRegistration RestaurantRegistration
	err := json.NewDecoder(r.Body).Decode(&restaurantRegistration)
	if err != nil {
		return nil, err
	}
	// XXX please look away, this is a hack
	if strings.Contains(strings.ToLower(restaurantRegistration.Restaurant.Address), "tartu") {
		restaurantRegistration.Region = "Tartu"
	} else {
		restaurantRegistration.Region = "Tallinn"
		// yup ...
	}
	return &restaurantRegistration, nil
}
