package main

import (
	"log"
	"net/http"

	"github.com/deiwin/praad-api/db"
	"github.com/deiwin/praad-api/handler"
	"github.com/gorilla/mux"
)

func main() {
	dbConfig := db.NewConfig()
	dbClient := db.NewClient(dbConfig)
	err := dbClient.Connect()
	if err != nil {
		panic(err)
	}
	defer dbClient.Disconnect()

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.HandleFunc("/offers", handler.Offers(dbClient))

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
