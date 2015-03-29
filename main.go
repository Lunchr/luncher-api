package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/deiwin/facebook"
	"github.com/deiwin/imstor"
	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/handler"
	"github.com/deiwin/luncher-api/router"
	"github.com/deiwin/luncher-api/session"
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

	imageSizes := []imstor.Size{
		imstor.Size{
			Name:   "large",
			Width:  800,
			Height: 400,
		},
	}
	imageFormats := []imstor.Format{
		imstor.PNG2JPEG,
		imstor.JPEGFormat,
	}
	imageStorageConf := imstor.NewConfig(imageSizes, imageFormats)
	imageStorage := imstor.New(imageStorageConf)

	r := router.NewWithPrefix("/api/v1/")
	r.GET("/offers", handler.Offers(offersCollection, regionsCollection, imageStorage))
	r.POST("/offers", handler.PostOffers(offersCollection, usersCollection, restaurantsCollection, sessionManager, facebookAuthenticator, imageStorage))
	r.GET("/tags", handler.Tags(tagsCollection))
	r.GET("/login/facebook", facebookHandler.Login())
	r.GET("/login/facebook/redirected", facebookHandler.Redirected())

	http.Handle("/api/v1/", r)
	portString := fmt.Sprintf(":%d", mainConfig.Port)
	log.Fatal(http.ListenAndServe(portString, nil))
}
