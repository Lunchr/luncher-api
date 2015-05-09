package handler

import (
	"net/http"
	"time"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"github.com/deiwin/luncher-api/router"
	"github.com/deiwin/luncher-api/session"
	"github.com/deiwin/luncher-api/storage"
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
