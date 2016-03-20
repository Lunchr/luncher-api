package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"github.com/Lunchr/luncher-api/storage"
	"github.com/deiwin/facebook"
	"github.com/julienschmidt/httprouter"
)

// UserRestaurants returns a list of restaurants the user has access to
func UserRestaurants(restaurants db.Restaurants, sessionManager session.Manager, users db.Users) router.Handler {
	handlerWithUser := func(w http.ResponseWriter, r *http.Request, user *model.User) *router.HandlerError {
		restaurantsByIDs, err := restaurants.GetByIDs(user.RestaurantIDs)
		if err != nil {
			return router.NewHandlerError(err, "Failed to find restaurants associated with this user", http.StatusInternalServerError)
		}
		fbPageAccessTokens := user.Session.FacebookPageTokens
		fbPageIDs := make([]string, len(fbPageAccessTokens))
		for i, fbPageAccessToken := range fbPageAccessTokens {
			fbPageIDs[i] = fbPageAccessToken.PageID
		}
		restaurantsByFBPageIDs, err := restaurants.GetByFacebookPageIDs(fbPageIDs)
		if err != nil {
			return router.NewHandlerError(err, "Failed to find restaurants for FB pages associated with this user", http.StatusInternalServerError)
		}
		allRestaurants := append(restaurantsByIDs, restaurantsByFBPageIDs...)
		return writeJSON(w, allRestaurants)
	}
	return checkLogin(sessionManager, users, handlerWithUser)
}

// Restaurant returns a router.Handler that returns the restaurant information for the specified restaurant
func Restaurant(c db.Restaurants, sessionManager session.Manager, users db.Users) router.HandlerWithParams {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User, restaurant *model.Restaurant) *router.HandlerError {
		return writeJSON(w, restaurant)
	}
	return forRestaurant(sessionManager, users, c, handler)
}

// PostRestaurants returns an handler for creating a restaurant
func PostRestaurants(c db.Restaurants, sessionManager session.Manager, users db.Users, fbAuth facebook.Authenticator) router.Handler {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User) *router.HandlerError {
		restaurant, err := parseRestaurant(r)
		if err != nil {
			return router.NewHandlerError(err, "Failed to parse the restaurant", http.StatusBadRequest)
		} else if restaurant.FacebookPageID == "" {
			return router.NewHandlerError(err, "Registering without an assciated FB page is currently disabled", http.StatusBadRequest)
		}
		insertedRestaurants, err := c.Insert(restaurant)
		if err != nil {
			return router.NewHandlerError(err, "Failed to store the restaurant in the DB", http.StatusInternalServerError)
		}
		var insertedRestaurant = insertedRestaurants[0]
		// We want to leave the FB page related restaurant role management wholly to FB, so we handle them totally
		// separately
		if insertedRestaurant.FacebookPageID == "" {
			user.RestaurantIDs = append(user.RestaurantIDs, insertedRestaurant.ID)
			err = users.Update(user.FacebookUserID, user)
			if err != nil {
				// TODO: revert the restaurant insertion we just did? Look into mgo's txn package
				return router.NewHandlerError(err, "Failed to attach the restaurant to the user in the DB", http.StatusInternalServerError)
			}
		} else {
			pageAccessToken, handlerErr := getPageAccessToken(&user.Session.FacebookUserToken, insertedRestaurant.FacebookPageID, fbAuth)
			if handlerErr != nil {
				return handlerErr
			}
			user.Session.FacebookPageTokens = append(user.Session.FacebookPageTokens, pageAccessToken)
			err = users.Update(user.FacebookUserID, user)
			if err != nil {
				return router.NewHandlerError(err, "Failed to store the access token for this restaurant's page in the DB", http.StatusInternalServerError)
			}
		}
		return writeJSON(w, insertedRestaurant)
	}
	return checkLogin(sessionManager, users, handler)
}

// RestaurantOffers returns all upcoming offers for the restaurant linked to the currently
// logged in user unless the request includes a 'title' query parameter, in which the offer
// with the specified title will be fetched instead.
func RestaurantOffers(restaurants db.Restaurants, sessionManager session.Manager, users db.Users, offers db.Offers,
	imageStorage storage.Images, regions db.Regions) router.HandlerWithParams {
	getTodaysOffersForRestaurant := func(w http.ResponseWriter, restaurant *model.Restaurant) *router.HandlerError {
		region, err := regions.GetName(restaurant.Region)
		if err != nil {
			return router.NewHandlerError(err, "Failed to find the region for this restaurant", http.StatusInternalServerError)
		}
		timeLocation, err := time.LoadLocation(region.Location)
		if err != nil {
			return router.NewHandlerError(err, "The location of this region is misconfigured", http.StatusInternalServerError)
		}
		today, _ := getTodaysTimeRange(timeLocation)
		offers, err := offers.GetForRestaurant(restaurant.ID, today)
		if err != nil {
			return router.NewHandlerError(err, "Failed to find upcoming offers for this restaurant", http.StatusInternalServerError)
		}
		offerJSONs, handlerErr := mapOffersToJSON(offers, imageStorage)
		if handlerErr != nil {
			return handlerErr
		}
		return writeJSON(w, offerJSONs)
	}

	getOfferByTitle := func(w http.ResponseWriter, restaurant *model.Restaurant, escapedTitle string) *router.HandlerError {
		title, err := url.QueryUnescape(escapedTitle)
		if err != nil {
			return router.NewHandlerError(err, "Failed to parse the specified title", http.StatusBadRequest)
		}
		offer, err := offers.GetForRestaurantByTitle(restaurant.ID, title)
		if err == mgo.ErrNotFound {
			return router.NewHandlerError(err, "Failed to find an offer with the specified title", http.StatusNotFound)
		} else if err != nil {
			return router.NewHandlerError(err, "Failed to find an offer with the specified title", http.StatusInternalServerError)
		}
		offerJSON, handlerErr := mapOfferToJSON(offer, imageStorage)
		if handlerErr != nil {
			return handlerErr
		}
		return writeJSON(w, offerJSON)
	}

	handler := func(w http.ResponseWriter, r *http.Request, user *model.User, restaurant *model.Restaurant) *router.HandlerError {
		title := r.FormValue("title")
		if title != "" {
			return getOfferByTitle(w, restaurant, title)
		}
		return getTodaysOffersForRestaurant(w, restaurant)
	}
	return forRestaurant(sessionManager, users, restaurants, handler)
}

// RestaurantOfferSuggestions handles GET requests to /restaurants/:id/offer_suggestions and expects a 'title' query
// parameter. It returns a list of previously used offer titles matching the one provided.
func RestaurantOfferSuggestions(restaurants db.Restaurants, sessionManager session.Manager, users db.Users,
	offers db.Offers) router.HandlerWithParams {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User, restaurant *model.Restaurant) *router.HandlerError {
		partialTitle := r.FormValue("title")
		if partialTitle == "" {
			return router.NewStringHandlerError("Title not specified!", "Please specify a title", http.StatusBadRequest)
		}
		matchingTitles, err := offers.GetSimilarTitlesForRestaurant(restaurant.ID, partialTitle)
		if err != nil {
			return router.NewHandlerError(err, "Failed to find matching offers", http.StatusInternalServerError)
		}
		return writeJSON(w, matchingTitles)
	}
	return forRestaurant(sessionManager, users, restaurants, handler)
}

type HandlerWithRestaurant func(w http.ResponseWriter, r *http.Request, user *model.User,
	restaurant *model.Restaurant) *router.HandlerError

func forRestaurant(sessionManager session.Manager, users db.Users, restaurants db.Restaurants, handler HandlerWithRestaurant) router.HandlerWithParams {
	handlerWithUser := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user *model.User) *router.HandlerError {
		restaurant, handlerErr := getRestaurantByParams(ps, user, restaurants)
		if handlerErr != nil {
			return handlerErr
		}
		return handler(w, r, user, restaurant)
	}
	return checkLoginWithParams(sessionManager, users, handlerWithUser)
}

type HandlerWithParamsWithRestaurant func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user *model.User,
	restaurant *model.Restaurant) *router.HandlerError

func forRestaurantWithParams(sessionManager session.Manager, users db.Users, restaurants db.Restaurants,
	handler HandlerWithParamsWithRestaurant) router.HandlerWithParams {
	handlerWithUser := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user *model.User) *router.HandlerError {
		restaurant, handlerErr := getRestaurantByParams(ps, user, restaurants)
		if handlerErr != nil {
			return handlerErr
		}
		return handler(w, r, ps, user, restaurant)
	}
	return checkLoginWithParams(sessionManager, users, handlerWithUser)
}

func getRestaurantByParams(ps httprouter.Params, user *model.User, restaurants db.Restaurants) (*model.Restaurant, *router.HandlerError) {
	restaurantIDString := ps.ByName("restaurantID")
	if restaurantIDString == "" {
		return nil, router.NewSimpleHandlerError("Expected a restaurant ID to be specified", http.StatusBadRequest)
	} else if !bson.IsObjectIdHex(restaurantIDString) {
		return nil, router.NewSimpleHandlerError("Invalid restaurant ID", http.StatusBadRequest)
	}
	restaurantID := bson.ObjectIdHex(restaurantIDString)
	restaurant, err := restaurants.GetID(restaurantID)
	if err == mgo.ErrNotFound {
		return nil, router.NewHandlerError(err, "Failed to find the specified restaurant", http.StatusNotFound)
	} else if err != nil {
		return nil, router.NewHandlerError(err, "Something went wrong while trying to find the specified restaurant", http.StatusInternalServerError)
	}
	if !authorizedToManageRestaurant(user, restaurant) {
		return nil, router.NewSimpleHandlerError("Not authorized to access this restaurant", http.StatusForbidden)
	}
	return restaurant, nil
}

func authorizedToManageRestaurant(user *model.User, restaurant *model.Restaurant) bool {
	hasDirectAccess := idsInclude(user.RestaurantIDs, restaurant.ID)
	authorizedThroughFB := hasPageAccessTokenForRestaurant(user, restaurant)
	return hasDirectAccess || authorizedThroughFB
}

func idsInclude(ids []bson.ObjectId, id bson.ObjectId) bool {
	for i := range ids {
		if ids[i] == id {
			return true
		}
	}
	return false
}

func hasPageAccessTokenForRestaurant(user *model.User, restaurant *model.Restaurant) bool {
	for _, pageAccessToken := range user.Session.FacebookPageTokens {
		if pageAccessToken.PageID == restaurant.FacebookPageID {
			return true
		}
	}
	return false
}

func parseRestaurant(r *http.Request) (*model.Restaurant, error) {
	var restaurant model.Restaurant
	err := json.NewDecoder(r.Body).Decode(&restaurant)
	if err != nil {
		return nil, err
	}
	// Add default values for configurable fields
	if restaurant.DefaultGroupPostMessageTemplate == "" {
		restaurant.DefaultGroupPostMessageTemplate = "Tänased päevapakkumised on:"
	}
	// XXX please look away, this is a hack
	if strings.Contains(strings.ToLower(restaurant.Address), "tartu") {
		restaurant.Region = "Tartu"
	} else {
		restaurant.Region = "Tallinn"
		// yup ...
	}
	return &restaurant, nil
}
