package handler

import (
	"log"
	"net/http"

	"github.com/deiwin/luncher-api/db"
)

func Tags(tagsCollection db.Tags) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tags, err := tagsCollection.Get()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		writeJson(w, tags)
	}
}
