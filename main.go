package main

import (
	"fmt"
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

	usersCollection := db.NewUsers(dbClient)
	offersCollection := db.NewOffers(dbClient)
	tagsCollection := db.NewTags(dbClient)

	sessionManager := session.NewManager()
	mainConfig, err := NewConfig()
	if err != nil {
		panic(err)
	}

	facebookConfig := facebook.NewConfig()
	facebookAuthenticator := facebook.NewAuthenticator(facebookConfig, mainConfig.Domain)
	facebookAPI := facebook.NewAPI(facebookConfig)
	facebookHandler := handler.NewFacebook(facebookAuthenticator, sessionManager, facebookAPI, usersCollection)

	r := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	r.Handle("/offers", handler.Offers(offersCollection))
	r.Handle("/tags", handler.Tags(tagsCollection))
	r.Handle("/login/facebook", facebookHandler.Login())
	r.Handle("/login/facebook/redirected", facebookHandler.Redirected())
	http.Handle("/", r)
	portString := fmt.Sprintf(":%d", mainConfig.Port)
	log.Fatal(http.ListenAndServe(portString, nil))
}
