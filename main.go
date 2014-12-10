package main

import (
	"log"
	"net/http"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/facebook"
	"github.com/deiwin/luncher-api/handler"
	"github.com/deiwin/luncher-api/session"
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

	facebookConfig := facebook.NewConfig()
	facebookAuthenticator, err := facebook.NewAuthenticator(facebookConfig, "localhost:8080")
	if err != nil {
		log.Fatal("Failed to set up facebook authenticator", err)
	}
	sessionManager := session.NewManager()

	offersCollection := db.NewOffers(dbClient)
	tagsCollection := db.NewTags(dbClient)

	r := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	r.HandleFunc("/offers", handler.Offers(offersCollection))
	r.HandleFunc("/tags", handler.Tags(tagsCollection))
	r.HandleFunc("/login/facebook", handler.FacebookLogin(facebookAuthenticator, sessionManager))
	r.HandleFunc("/login/facebook", handler.FacebookRedirected(facebookAuthenticator, sessionManager))
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
