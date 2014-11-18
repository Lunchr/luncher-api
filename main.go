package main

import (
	"log"
	"net/http"

	"github.com/deiwin/praad-api/db"
	"github.com/deiwin/praad-api/handler"
	"github.com/gorilla/mux"
)

func main() {
	dbClient := db.NewClient()
	err := dbClient.Connect()
	if err != nil {
		panic(err)
	}
	defer dbClient.Disconnect()

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.HandleFunc("/offers", handler.Offers)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
