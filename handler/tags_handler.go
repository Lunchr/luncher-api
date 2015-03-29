package handler

import (
	"net/http"

	"github.com/deiwin/luncher-api/db"
	. "github.com/deiwin/luncher-api/router"
)

func Tags(tagsCollection db.Tags) Handler {
	return func(w http.ResponseWriter, r *http.Request) *HandlerError {
		tags, err := tagsCollection.Get()
		if err != nil {
			return &HandlerError{err, "", http.StatusInternalServerError}
		}
		return writeJSON(w, tags)
	}
}
