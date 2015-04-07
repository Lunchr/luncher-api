package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/deiwin/facebook"
	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/handler"
	"github.com/deiwin/luncher-api/router"
	"github.com/deiwin/luncher-api/session"
	"github.com/deiwin/luncher-api/storage"
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
	offersCollection, err := db.NewOffers(dbClient)
	if err != nil {
		panic(err)
	}
	tagsCollection := db.NewTags(dbClient)
	regionsCollection := db.NewRegions(dbClient)
	restaurantsCollection := db.NewRestaurants(dbClient)

	sessionManager := session.NewManager()
	mainConfig, err := NewConfig()
	if err != nil {
		panic(err)
	}

	redirectURL := mainConfig.Domain + "/api/v1/login/facebook/redirected"
	scopes := []string{"manage_pages", "publish_actions"}
	facebookConfig := facebook.NewConfig(redirectURL, scopes)
	facebookAuthenticator := facebook.NewAuthenticator(facebookConfig)
	facebookHandler := handler.NewFacebook(facebookAuthenticator, sessionManager, usersCollection)

	imageStorage := storage.NewImages()

	r := router.NewWithPrefix("/api/v1/")
	r.GET("/offers", handler.Offers(offersCollection, regionsCollection, imageStorage))
	r.POST("/offers", handler.PostOffers(offersCollection, usersCollection, restaurantsCollection, sessionManager, facebookAuthenticator, imageStorage))
	r.PUTWithParams("/offers/:id", handler.PutOffers(offersCollection, usersCollection, restaurantsCollection, sessionManager, facebookAuthenticator, imageStorage))
	r.GET("/tags", handler.Tags(tagsCollection))
	r.GET("/restaurant", handler.Restaurant(restaurantsCollection, sessionManager, usersCollection))
	r.GET("/restaurant/offers", handler.RestaurantOffers(restaurantsCollection, sessionManager, usersCollection, offersCollection, imageStorage))
	r.GET("/login/facebook", facebookHandler.Login())
	r.GET("/login/facebook/redirected", facebookHandler.Redirected())

	http.Handle("/api/v1/", r)
	portString := fmt.Sprintf(":%d", mainConfig.Port)
	log.Fatal(http.ListenAndServe(portString, nil))
}
