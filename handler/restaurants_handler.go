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
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User, restaurant *model.Restaurant) *router.HandlerError {
		return writeJSON(w, restaurant)
	}
	return forRestaurant(sessionManager, users, c, handler)
}

// Restaurant returns a router.Handler that returns the restaurant information for the
// restaurant linked to the currently logged in user
func PostRestaurants(c db.Restaurants, sessionManager session.Manager, users db.Users) router.Handler {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User) *router.HandlerError {
		restaurant, err := parseRestaurant(r)
		if err != nil {
			return router.NewHandlerError(err, "Failed to parse the restaurant", http.StatusBadRequest)
		}
		insertedRestaurants, err := c.Insert(restaurant)
		if err != nil {
			return router.NewHandlerError(err, "Failed to store the restaurant in the DB", http.StatusInternalServerError)
		}
		var insertedRestaurant = insertedRestaurants[0]
		user.RestaurantIDs = append(user.RestaurantIDs, insertedRestaurant.ID)
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
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User, restaurant *model.Restaurant) *router.HandlerError {
		offers, err := offers.GetForRestaurant(restaurant.Name, time.Now())
		if err != nil {
			return router.NewHandlerError(err, "Failed to find upcoming offers for this restaurant", http.StatusInternalServerError)
		}
		offerJSONs, handlerErr := mapOffersToJSON(offers, imageStorage)
		if handlerErr != nil {
			return handlerErr
		}
		return writeJSON(w, offerJSONs)
	}
	return forRestaurant(sessionManager, users, restaurants, handler)
}

type HandlerWithRestaurant func(w http.ResponseWriter, r *http.Request, user *model.User, restaurant *model.Restaurant) *router.HandlerError

func forRestaurant(sessionManager session.Manager, users db.Users, restaurants db.Restaurants, handler HandlerWithRestaurant) router.Handler {
	handlerWithUser := func(w http.ResponseWriter, r *http.Request, user *model.User) *router.HandlerError {
		restaurant, err := restaurants.GetID(user.RestaurantIDs[0])
		if err != nil {
			return router.NewHandlerError(err, "Failed to find the restaurant connected to this user", http.StatusInternalServerError)
		}
		return handler(w, r, user, restaurant)
	}
	return checkLogin(sessionManager, users, handlerWithUser)
}

func parseRestaurant(r *http.Request) (*model.Restaurant, error) {
	var restaurant model.Restaurant
	err := json.NewDecoder(r.Body).Decode(&restaurant)
	if err != nil {
		return nil, err
	}
	// XXX please look away, this is a hack
	if strings.Contains(strings.ToLower(restaurant.Address), "tartu") {
		restaurant.Region = "Tartu"
	} else {
		restaurant.Region = "Tallinn"
		// yup ...
	}
	return &restaurant, nil
}
