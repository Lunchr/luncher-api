package handler

import (
	"net/http"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/router"
)

func Tags(tagsCollection db.Tags) router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		tags, err := tagsCollection.Get()
		if err != nil {
			return &router.HandlerError{err, "", http.StatusInternalServerError}
		}
		return writeJSON(w, tags)
	}
}
