package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/deiwin/facebook"
	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"github.com/deiwin/luncher-api/session"
)

// Offers handles GET requests to /offers. It returns all current day's offers.
func Offers(offersCollection db.Offers) Handler {
	return func(w http.ResponseWriter, r *http.Request) *handlerError {
		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		endTime := startTime.AddDate(0, 0, 1)
		offers, err := offersCollection.GetForTimeRange(startTime, endTime)
		if err != nil {
			return &handlerError{err, "", http.StatusInternalServerError}
		}
		return writeJSON(w, offers)
	}
}

// PostOffers handles POST requests to /offers. It stores the offer in the DB and
// sends it to Facebook to be posted on the page's wall at the requested time.
func PostOffers(offersCollection db.Offers, usersCollection db.Users, sessionManager session.Manager, fbAuth facebook.Authenticator) Handler {
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
		offer := parseOffer(r)
		// if err != nil {
		// 	return &handlerError{err, "", http.StatusBadRequest}
		// } TODO maybe check that all the required fields are actually set?
		message := formFBOfferMessage(offer)
		post, err := api.PagePublish(user.Session.FacebookPageToken, user.FacebookPageID, message)
		if err != nil {
			return &handlerError{err, "", http.StatusBadGateway}
		}
		offer.FBPostID = post.ID
		if err := offersCollection.Insert(&offer); err != nil {
			return &handlerError{err, "", http.StatusInternalServerError}
		}

		return nil
	}
}

func formFBOfferMessage(o model.Offer) string {
	return fmt.Sprintf("%s - %s", o.Title, o.Description)
}

func parseOffer(r *http.Request) model.Offer {
	return model.Offer{
		Title:       r.PostFormValue("title"),
		Description: r.PostFormValue("description"),
	}
}
