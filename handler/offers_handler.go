package handler

import (
	"net/http"
	"time"

	"github.com/deiwin/facebook"
	"github.com/deiwin/luncher-api/db"
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
		return nil
	}
}
