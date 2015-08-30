package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"github.com/deiwin/facebook"
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
func PostOfferGroupPost(c db.OfferGroupPosts, sessionManager session.Manager, users db.Users, restaurants db.Restaurants,
	offers db.Offers, regions db.Regions, fbAuth facebook.Authenticator) router.Handler {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User, restaurant *model.Restaurant) *router.HandlerError {
		post, handlerErr := parseOfferGroupPost(r, restaurant)
		if handlerErr != nil {
			return handlerErr
		}
		insertedPosts, err := c.Insert(post)
		if err != nil {
			return router.NewHandlerError(err, "Failed to store the post in the DB", http.StatusInternalServerError)
		}
		insertedPost := insertedPosts[0]
		if handlerErr = updateGroupPost(insertedPost, user, restaurant, offers, regions, c, fbAuth); handlerErr != nil {
			return handlerErr
		}
		return writeJSON(w, insertedPost)
	}
	return forRestaurant(sessionManager, users, restaurants, handler)
}

// PutOfferGroupPost handles PUT requests to /restaurant/posts/:date. It stores the info in the DB and updates the post in FB.
func PutOfferGroupPost(c db.OfferGroupPosts, sessionManager session.Manager, users db.Users, restaurants db.Restaurants,
	offers db.Offers, regions db.Regions, fbAuth facebook.Authenticator) router.HandlerWithParams {
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
		// Update by date, because the posted data does not include the previous FB post ID
		if handlerErr = updateGroupPostForDate(post.Date, user, restaurant, offers, regions, c, fbAuth); handlerErr != nil {
			return handlerErr
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

func updateGroupPostForDate(date model.DateWithoutTime, user *model.User, restaurant *model.Restaurant, offers db.Offers,
	regions db.Regions, groupPosts db.OfferGroupPosts, fbAuth facebook.Authenticator) *router.HandlerError {
	post, err := groupPosts.GetByDate(date, restaurant.ID)
	if err == mgo.ErrNotFound {
		postToInsert := &model.OfferGroupPost{
			RestaurantID:    restaurant.ID,
			Date:            date,
			MessageTemplate: restaurant.DefaultGroupPostMessageTemplate,
		}
		insertedPosts, err := groupPosts.Insert(postToInsert)
		if err != nil {
			router.NewHandlerError(err, "Failed to create a group post with restaurant defaults", http.StatusInternalServerError)
		}
		post = insertedPosts[0]
	} else if err != nil {
		return router.NewHandlerError(err, "Failed to fetch a group post for that date", http.StatusInternalServerError)
	}
	return updateGroupPost(post, user, restaurant, offers, regions, groupPosts, fbAuth)
}

func updateGroupPost(post *model.OfferGroupPost, user *model.User, restaurant *model.Restaurant, offers db.Offers,
	regions db.Regions, groupPosts db.OfferGroupPosts, fbAuth facebook.Authenticator) *router.HandlerError {
	if restaurant.FacebookPageID == "" {
		return nil
	}
	fbAPI := fbAuth.APIConnection(&user.Session.FacebookUserToken)
	// Remove the current post from FB, if it's already there
	if post.FBPostID != "" {
		err := fbAPI.PostDelete(user.Session.FacebookPageToken, post.FBPostID)
		if err != nil {
			return router.NewHandlerError(err, "Failed to delete the current post from Facebook", http.StatusBadGateway)
		}
	}
	offersForDate, handlerErr := getOffersForDate(post.Date, restaurant, offers, regions)
	if handlerErr != nil {
		return handlerErr
	} else if len(offersForDate) == 0 {
		return nil
	}
	message := formFBMessage(post, offersForDate)
	// Add the new version
	fbPost, err := fbAPI.PagePublish(user.Session.FacebookPageToken, restaurant.FacebookPageID, message)
	if err != nil {
		return router.NewHandlerError(err, "Failed to post the offer to Facebook", http.StatusBadGateway)
	}
	post.FBPostID = fbPost.ID
	if err = groupPosts.UpdateByID(post.ID, post); err != nil {
		return router.NewHandlerError(err, "Failed to update a group post in the DB", http.StatusInternalServerError)
	}
	return nil
}

func formFBMessage(post *model.OfferGroupPost, offers []*model.Offer) string {
	offerMessages := make([]string, len(offers))
	for i, offer := range offers {
		offerMessages[i] = formFBOfferMessage(offer)
	}
	offersMessage := strings.Join(offerMessages, "\n")
	return fmt.Sprintf("%s\n\n%s", post.MessageTemplate, offersMessage)
}

func formFBOfferMessage(o *model.Offer) string {
	// TODO get rid of the hard-coded €
	return fmt.Sprintf("%s - %.2f€", o.Title, o.Price)
}

func getOffersForDate(date model.DateWithoutTime, restaurant *model.Restaurant, offers db.Offers, regions db.Regions) ([]*model.Offer, *router.HandlerError) {
	region, err := regions.GetName(restaurant.Region)
	if err != nil {
		return nil, router.NewHandlerError(err, "Failed to find the restaurant's region", http.StatusInternalServerError)
	}
	location, err := time.LoadLocation(region.Location)
	if err != nil {
		return nil, router.NewHandlerError(err, "Failed to load region's location", http.StatusInternalServerError)
	}
	startTime, endTime, err := date.TimeBounds(location)
	if err != nil {
		return nil, router.NewHandlerError(err, "Failed to parse a date", http.StatusInternalServerError)
	}
	offersForDate, err := offers.GetForRestaurantWithinTimeBounds(restaurant.ID, startTime, endTime)
	if err != nil {
		return nil, router.NewHandlerError(err, "Failed to find offers for this date", http.StatusInternalServerError)
	}
	return offersForDate, nil
}
