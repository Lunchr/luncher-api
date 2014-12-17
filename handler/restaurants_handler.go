package handler

import (
	"log"
	"net/http"

	"github.com/deiwin/luncher-api/db"
)

func Restaurants(restaurantsCollection db.Restaurants) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		restaurants, err := restaurantsCollection.Get()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		writeJSON(w, restaurants)
	}
}
