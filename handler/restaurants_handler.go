package handler

import (
	"net/http"

	"github.com/deiwin/luncher-api/db"
	. "github.com/deiwin/luncher-api/router"
)

func Restaurants(restaurantsCollection db.Restaurants) Handler {
	return func(w http.ResponseWriter, r *http.Request) *HandlerError {
		restaurants, err := restaurantsCollection.Get()
		if err != nil {
			return &HandlerError{err, "", http.StatusInternalServerError}
		}
		return writeJSON(w, restaurants)
	}
}
