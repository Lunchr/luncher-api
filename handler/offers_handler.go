package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/deiwin/facebook"
	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"github.com/deiwin/luncher-api/session"
)

// Offers handles GET requests to /offers. It returns all current day's offers.
func Offers(offersCollection db.Offers, regionsCollection db.Regions) Handler {
	return func(w http.ResponseWriter, r *http.Request) *handlerError {
		regionName := r.FormValue("region")
		if regionName == "" {
			return &handlerError{errors.New("Region not specified for GET /offers"), "Please specify a region", http.StatusBadRequest}
		}
		region, err := regionsCollection.Get(regionName)
		if err != nil {
			return &handlerError{err, "", http.StatusInternalServerError}
		}
		loc, err := time.LoadLocation(region.Location)
		if err != nil {
			return &handlerError{err, "", http.StatusInternalServerError}
		}
		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		endTime := startTime.AddDate(0, 0, 1)
		offers, err := offersCollection.Get(regionName, startTime, endTime)
		if err != nil {
			return &handlerError{err, "", http.StatusInternalServerError}
		}
		return writeJSON(w, offers)
	}
}

// PostOffers handles POST requests to /offers. It stores the offer in the DB and
// sends it to Facebook to be posted on the page's wall at the requested time.
func PostOffers(offersCollection db.Offers, usersCollection db.Users, restaurantsCollection db.Restaurants,
	sessionManager session.Manager, fbAuth facebook.Authenticator) Handler {
	return func(w http.ResponseWriter, r *http.Request) *handlerError {
		session, err := sessionManager.Get(r)
		if err != nil {
			return &handlerError{err, "", http.StatusForbidden}
		}
		user, err := usersCollection.GetBySessionID(session)
		if err != nil {
			return &handlerError{err, "", http.StatusForbidden}
		}
		api := fbAuth.APIConnection(&user.Session.FacebookUserToken)
		restaurant, err := restaurantsCollection.GetByID(user.RestaurantID)
		if err != nil {
			return &handlerError{err, "", http.StatusInternalServerError}
		}
		offer, err := parseOffer(r, restaurant)
		if err != nil {
			return &handlerError{err, "", http.StatusBadRequest}
		}
		message := formFBOfferMessage(*offer)
		post, err := api.PagePublish(user.Session.FacebookPageToken, user.FacebookPageID, message)
		if err != nil {
			return &handlerError{err, "", http.StatusBadGateway}
		}
		offer.FBPostID = post.ID
		offers, err := offersCollection.Insert(offer)
		if err != nil {
			return &handlerError{err, "", http.StatusInternalServerError}
		}

		return writeJSON(w, offers[0])
	}
}

func formFBOfferMessage(o model.Offer) string {
	return fmt.Sprintf("%s - %s", o.Title, o.Description)
}

func parseOffer(r *http.Request, restaurant *model.Restaurant) (*model.Offer, error) {
	price, err := strconv.ParseFloat(r.PostFormValue("price"), 64)
	if err != nil {
		return nil, err
	}

	fromTime, err := time.Parse(time.RFC3339, r.PostFormValue("from_time"))
	if err != nil {
		return nil, err
	}
	toTime, err := time.Parse(time.RFC3339, r.PostFormValue("to_time"))
	if err != nil {
		return nil, err
	}

	offer := &model.Offer{
		Title:       r.PostFormValue("title"),
		Description: r.PostFormValue("description"),
		Tags:        r.Form["tags"],
		Price:       price,
		Restaurant: model.OfferRestaurant{
			Name:   restaurant.Name,
			Region: restaurant.Region,
		},
		FromTime: fromTime,
		ToTime:   toTime,
	}
	return offer, nil
}
