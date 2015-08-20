package handler

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"github.com/julienschmidt/httprouter"
)

// OfferGroupPost handles GET requests to /restaurant/posts/:date. It returns all current day's offers for the region.
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

// PostOfferGroupPost handles POST requests to /restaurant/posts. It stores the info in the DB and updates the post in FB.
func PostOfferGroupPost(c db.OfferGroupPosts, sessionManager session.Manager, users db.Users, restaurants db.Restaurants) router.Handler {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User, restaurant *model.Restaurant) *router.HandlerError {
		post, handlerErr := parseOfferGroupPost(r, restaurant)
		if handlerErr != nil {
			return handlerErr
		}
		insertedPosts, err := c.Insert(post)
		if err != nil {
			return router.NewHandlerError(err, "Failed to store the post in the DB", http.StatusInternalServerError)
		}
		return writeJSON(w, insertedPosts[0])
	}
	return forRestaurant(sessionManager, users, restaurants, handler)
}

// PutOfferGroupPost handles PUT requests to /restaurant/posts/:date. It stores the info in the DB and updates the post in FB.
func PutOfferGroupPost(c db.OfferGroupPosts, sessionManager session.Manager, users db.Users, restaurants db.Restaurants) router.HandlerWithParams {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User, restaurant *model.Restaurant,
		date model.DateWithoutTime) *router.HandlerError {
		post, handlerErr := parseOfferGroupPost(r, restaurant)
		if handlerErr != nil {
			return handlerErr
		}
		if post.Date != date {
			return router.NewSimpleHandlerError("Unexpected date value", http.StatusBadRequest)
		}
		err := c.UpdateByID(post.ID, post)
		if err != nil {
			return router.NewSimpleHandlerError("Failed to insert the post to DB", http.StatusBadRequest)
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

func parseOfferGroupPost(r *http.Request, restaurant *model.Restaurant) (*model.OfferGroupPost, *router.HandlerError) {
	var post struct {
		ID              bson.ObjectId `json:"_id"`
		MessageTemplate string        `json:"message_template"`
		Date            string        `json:"date"`
	}
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		return nil, router.NewHandlerError(err, "Failed to parse the post", http.StatusBadRequest)
	}
	date := model.DateWithoutTime(post.Date)
	if date == "" {
		return nil, router.NewStringHandlerError("Date not specified!", "Please specify a date", http.StatusBadRequest)
	}
	if !date.IsValid() {
		return nil, router.NewSimpleHandlerError("Invalid date specified", http.StatusBadRequest)
	}
	return &model.OfferGroupPost{
		ID:              post.ID,
		MessageTemplate: post.MessageTemplate,
		Date:            date,
		RestaurantID:    restaurant.ID,
	}, nil
}
