package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/deiwin/praad-api/db"
)

func Offers(dbClient *db.Client) func(http.ResponseWriter, *http.Request) {
	offersCollection := db.NewOffers(dbClient)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		endTime := startTime.AddDate(0, 0, 1)
		offers, err := offersCollection.GetForTimeRange(startTime, endTime)
		if err != nil {
			log.Println(err)
		}
		writeJson(w, offers)
	}
}
