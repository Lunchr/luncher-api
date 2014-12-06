package main

import (
	"log"
	"net/http"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/handler"
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

	offersCollection := db.NewOffers(dbClient)
	tagsCollection := db.NewTags(dbClient)

	r := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	r.HandleFunc("/offers", handler.Offers(offersCollection))
	r.HandleFunc("/tags", handler.Tags(tagsCollection))
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
