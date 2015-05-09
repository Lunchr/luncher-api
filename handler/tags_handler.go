package handler

import (
	"net/http"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"github.com/deiwin/luncher-api/router"
)

func Tags(tagsCollection db.Tags) router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		tagsIter := tagsCollection.GetAll()
		var tags []model.Tag
		var tag model.Tag
		for tagsIter.Next(&tag) {
			tags = append(tags, tag)
		}
		if err := tagsIter.Close(); err != nil {
			return router.NewHandlerError(err, "An error occured while fetching the tags from the DB", http.StatusInternalServerError)
		}
		return writeJSON(w, tags)
	}
}
