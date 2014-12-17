package handler

import (
	"log"
	"net/http"

	"github.com/deiwin/luncher-api/db"
)

func Tags(tagsCollection db.Tags) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		tags, err := tagsCollection.Get()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		writeJSON(w, tags)
	}
}
