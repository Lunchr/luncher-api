package handler

import (
	"net/http"

	"github.com/deiwin/luncher-api/db"
)

func Tags(tagsCollection db.Tags) Handler {
	return func(w http.ResponseWriter, r *http.Request) *handlerError {
		tags, err := tagsCollection.Get()
		if err != nil {
			return &handlerError{err, "", http.StatusInternalServerError}
		}
		return writeJSON(w, tags)
	}
}
