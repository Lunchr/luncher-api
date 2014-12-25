package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Handler func(http.ResponseWriter, *http.Request)

func writeJSON(w http.ResponseWriter, v interface{}) {
	if data, err := json.Marshal(v); err != nil {
		log.Printf("Error marshalling json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}
