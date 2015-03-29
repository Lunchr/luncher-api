package handler

import (
	"net/http"
	"time"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	. "github.com/deiwin/luncher-api/router"
	"github.com/deiwin/luncher-api/session"
)

// Restaurants returns a list of all restaurants
func Restaurants(restaurantsCollection db.Restaurants) Handler {
	return func(w http.ResponseWriter, r *http.Request) *HandlerError {
		restaurants, err := restaurantsCollection.Get()
		if err != nil {
			return &HandlerError{err, "", http.StatusInternalServerError}
		}
		return writeJSON(w, restaurants)
	}
}

// Restaurant returns a Handler that returns the restaurant information for the
// restaurant linked to the currently logged in user
func Restaurant(c db.Restaurants, sessionManager session.Manager, users db.Users) Handler {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User) *HandlerError {
		restaurant, err := c.GetByID(user.RestaurantID)
		if err != nil {
			return &HandlerError{err, "", http.StatusInternalServerError}
		}
		return writeJSON(w, restaurant)
	}
	return checkLogin(sessionManager, users, handler)
}

// RestaurantOffers returns all upcoming offers for the restaurant linked to the
// currently logged in user
func RestaurantOffers(restaurants db.Restaurants, sessionManager session.Manager, users db.Users, offers db.Offers) Handler {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User) *HandlerError {
		restaurant, err := restaurants.GetByID(user.RestaurantID)
		if err != nil {
			return &HandlerError{err, "", http.StatusInternalServerError}
		}
		offers, err := offers.GetForRestaurant(restaurant.Name, time.Now())
		if err != nil {
			return &HandlerError{err, "", http.StatusInternalServerError}
		}
		return writeJSON(w, offers)
	}
	return checkLogin(sessionManager, users, handler)
}
