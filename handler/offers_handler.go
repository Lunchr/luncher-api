package handler

import (
	"net/http"
	"time"

	"github.com/deiwin/luncher-api/db"
)

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
