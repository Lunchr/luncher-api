package handler

import (
	"log"
	"net/http"

	"github.com/deiwin/praad-api/db"
)

func Restaurants(restaurantsCollection db.Restaurants) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		restaurants, err := restaurantsCollection.Get()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		writeJson(w, restaurants)
	}
}
