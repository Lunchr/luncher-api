package main

import (
	"github.com/deiwin/praad-api/handler"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.HandleFunc("/offers", handler.Offers)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
