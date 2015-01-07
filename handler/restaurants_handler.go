package handler

import (
	"net/http"

	"github.com/deiwin/luncher-api/db"
)

func Restaurants(restaurantsCollection db.Restaurants) Handler {
	return func(w http.ResponseWriter, r *http.Request) *handlerError {
		restaurants, err := restaurantsCollection.Get()
		if err != nil {
			return &handlerError{err, "", http.StatusInternalServerError}
		}
		return writeJSON(w, restaurants)
	}
}
