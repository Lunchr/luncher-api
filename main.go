package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/handler"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"github.com/Lunchr/luncher-api/storage"
	"github.com/deiwin/facebook"
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
	offerGroupPostsCollection := db.NewOfferGroupPosts(dbClient)

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
	r.POST("/restaurants", handler.PostRestaurants(restaurantsCollection, sessionManager, usersCollection))
	r.GET("/restaurant/offers", handler.RestaurantOffers(restaurantsCollection, sessionManager, usersCollection, offersCollection, imageStorage))
	r.GETWithParams("/restaurant/posts/:date", handler.OfferGroupPost(offerGroupPostsCollection, sessionManager, usersCollection, restaurantsCollection))
	r.POST("/restaurant/posts", handler.PostOfferGroupPost(offerGroupPostsCollection, sessionManager, usersCollection,
		restaurantsCollection, offersCollection, regionsCollection, facebookLoginAuthenticator))
	r.PUT("/restaurant/posts/:date", handler.PutOfferGroupPost(offerGroupPostsCollection, sessionManager, usersCollection,
		restaurantsCollection, offersCollection, regionsCollection, facebookLoginAuthenticator))
	r.GET("/logout", handler.Logout(sessionManager, usersCollection))
	r.GET("/login/facebook", handler.RedirectToFBForLogin(sessionManager, facebookLoginAuthenticator))
	r.GET("/login/facebook/redirected", handler.RedirectedFromFBForLogin(sessionManager, facebookLoginAuthenticator, usersCollection, restaurantsCollection))
	r.GET("/register/facebook", handler.RedirectToFBForRegistration(sessionManager, facebookRegistrationAuthenticator))
	r.GET("/register/facebook/redirected", handler.RedirectedFromFBForRegistration(sessionManager, facebookRegistrationAuthenticator, usersCollection))
	r.GET("/register/facebook/pages", handler.ListPagesManagedByUser(sessionManager, facebookRegistrationAuthenticator, usersCollection))
	r.GETWithParams("/register/facebook/pages/:id", handler.Page(sessionManager, facebookRegistrationAuthenticator, usersCollection))

	http.Handle("/api/v1/", r)
	portString := fmt.Sprintf(":%d", mainConfig.Port)
	log.Fatal(http.ListenAndServe(portString, nil))
}
