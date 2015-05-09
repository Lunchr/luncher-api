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

	scopes := []string{"manage_pages", "publish_pages"}
	loginRedirectURL := mainConfig.Domain + "/api/v1/login/facebook/redirected"
	facebookLoginConfig := facebook.NewConfig(loginRedirectURL, scopes)
	facebookLoginAuthenticator := facebook.NewAuthenticator(facebookLoginConfig)

	registrationRedirectURL := mainConfig.Domain + "/api/v1/register/facebook/redirected"
	facebookRegistrationConfig := facebook.NewConfig(registrationRedirectURL, scopes)
	facebookRegistrationAuthenticator := facebook.NewAuthenticator(facebookRegistrationConfig)

	facebookHandler := handler.NewFacebook(facebookLoginAuthenticator, facebookRegistrationAuthenticator, sessionManager, usersCollection)

	imageStorage := storage.NewImages()

	r := router.NewWithPrefix("/api/v1/")
	r.GET("/regions", handler.Regions(regionsCollection))
	r.GETWithParams("/regions/:name/offers", handler.RegionOffers(offersCollection, regionsCollection, imageStorage))
	r.GET("/offers", handler.ProximalOffers(offersCollection, imageStorage))
	r.POST("/offers", handler.PostOffers(offersCollection, usersCollection, restaurantsCollection, sessionManager, facebookLoginAuthenticator, imageStorage))
	r.PUT("/offers/:id", handler.PutOffers(offersCollection, usersCollection, restaurantsCollection, sessionManager, facebookLoginAuthenticator, imageStorage))
	r.DELETE("/offers/:id", handler.DeleteOffers(offersCollection, usersCollection, sessionManager, facebookLoginAuthenticator))
	r.GET("/tags", handler.Tags(tagsCollection))
	r.GET("/restaurant", handler.Restaurant(restaurantsCollection, sessionManager, usersCollection))
	r.GET("/restaurant/offers", handler.RestaurantOffers(restaurantsCollection, sessionManager, usersCollection, offersCollection, imageStorage))
	r.GET("/logout", handler.Logout(sessionManager, usersCollection))
	r.GET("/login/facebook", facebookHandler.RedirectToFBForLogin())
	r.GET("/login/facebook/redirected", facebookHandler.RedirectedFromFBForLogin())

	http.Handle("/api/v1/", r)
	portString := fmt.Sprintf(":%d", mainConfig.Port)
	log.Fatal(http.ListenAndServe(portString, nil))
}
