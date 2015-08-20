package handler

import (
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"github.com/julienschmidt/httprouter"
)

// OfferGroupPost handles GET requests to /restaurant/post/:date. It returns all current day's offers for the region.
func OfferGroupPost(c db.OfferGroupPosts, sessionManager session.Manager, users db.Users, restaurants db.Restaurants) router.HandlerWithParams {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User, restaurant *model.Restaurant,
		date model.DateWithoutTime) *router.HandlerError {
		post, err := c.GetByDate(date, restaurant.ID)
		if err == mgo.ErrNotFound {
			return router.NewHandlerError(err, "Offer group post not found", http.StatusNotFound)
		} else if err != nil {
			return router.NewHandlerError(err, "An error occured while trying to fetch a offer group post", http.StatusInternalServerError)
		}
		return writeJSON(w, post)
	}
	return forDate(sessionManager, users, restaurants, handler)
}

type HandlerWithRestaurantAndDate func(w http.ResponseWriter, r *http.Request, user *model.User, restaurant *model.Restaurant,
	date model.DateWithoutTime) *router.HandlerError

func forDate(sessionManager session.Manager, users db.Users, restaurants db.Restaurants,
	handler HandlerWithRestaurantAndDate) router.HandlerWithParams {
	handlerWithRestaurant := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user *model.User,
		restaurant *model.Restaurant) *router.HandlerError {
		date := model.DateWithoutTime(ps.ByName("date"))
		if date == "" {
			return router.NewStringHandlerError("Date not specified!", "Please specify a date", http.StatusBadRequest)
		}
		if !date.IsValid() {
			return router.NewSimpleHandlerError("Invalid date specified", http.StatusBadRequest)
		}
		return handler(w, r, user, restaurant, date)
	}
	return forRestaurantWithParams(sessionManager, users, restaurants, handlerWithRestaurant)
}
