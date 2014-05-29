package main

import (
	"github.com/deiwin/praad-api/db"
	"github.com/deiwin/praad-api/handler"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	db.Connect()
	defer db.Disconnect()

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.HandleFunc("/offers", handler.Offers)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
