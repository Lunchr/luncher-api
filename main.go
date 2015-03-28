package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/deiwin/facebook"
	"github.com/deiwin/imstor"
	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/handler"
	"github.com/deiwin/luncher-api/session"
	"github.com/julienschmidt/httprouter"
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

	r := httprouter.New()
	pathPrefix := "/api/v1"
	r.Handler("GET", pathPrefix+"/offers", handler.Offers(offersCollection, regionsCollection, imageStorage))
	r.Handler("POST", pathPrefix+"/offers", handler.PostOffers(offersCollection, usersCollection,
		restaurantsCollection, sessionManager, facebookAuthenticator, imageStorage))
	r.Handler("GET", pathPrefix+"/tags", handler.Tags(tagsCollection))
	r.Handler("GET", pathPrefix+"/login/facebook", facebookHandler.Login())
	r.Handler("GET", pathPrefix+"/login/facebook/redirected", facebookHandler.Redirected())
	http.Handle("/", r)
	portString := fmt.Sprintf(":%d", mainConfig.Port)
	log.Fatal(http.ListenAndServe(portString, nil))
}
