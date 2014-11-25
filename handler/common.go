package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func writeJson(w http.ResponseWriter, v interface{}) {
	if data, err := json.Marshal(v); err != nil {
		log.Printf("Error marshalling json: %v", err)
	} else {
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}