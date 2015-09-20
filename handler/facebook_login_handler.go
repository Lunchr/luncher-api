package handler

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"github.com/deiwin/facebook"
	"golang.org/x/oauth2"
)

// RedirectToFBForLogin returns a handler that redirects the user to Facebook to log in
func RedirectToFBForLogin(sessionManager session.Manager, auther facebook.Authenticator) router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		session := sessionManager.GetOrInit(w, r)
		redirectURL := auther.AuthURL(session)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return nil
	}
}

// RedirectedFromFBForLogin returns a handler that receives the user and page tokens for the
// user who has just logged in through Facebook. Updates the user and page
// access tokens in the DB
func RedirectedFromFBForLogin(sessionManager session.Manager, fbAuth facebook.Authenticator, users db.Users, restaurants db.Restaurants) router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		session := sessionManager.GetOrInit(w, r)
		tok, handlerErr := getLongTermToken(session, r, fbAuth)
		if handlerErr != nil {
			return handlerErr
		}
		fbUserID, err := getUserID(tok, fbAuth)
		if err != nil {
			return router.NewHandlerError(err, "Failed to get the user information from Facebook", http.StatusInternalServerError)
		}
		user, err := users.GetFbID(fbUserID)
		if err == mgo.ErrNotFound {
			return router.NewHandlerError(err, "User not registered", http.StatusForbidden)
		} else if err != nil {
			return router.NewHandlerError(err, "Failed to find the user from DB", http.StatusInternalServerError)
		}
		if handlerErr = storeAccessTokensInDB(user.ID, fbUserID, tok, session, users); err != nil {
			return handlerErr
		}
		if handlerErr = storeTokensForRestaurantPages(fbUserID, tok, restaurants, users, fbAuth); err != nil {
			return handlerErr
		}
		http.Redirect(w, r, "/#/admin", http.StatusSeeOther)
		return nil
	}
}

func storeTokensForRestaurantPages(fbUserID string, userAccessToken *oauth2.Token, restaurants db.Restaurants, users db.Users,
	fbAuth facebook.Authenticator) *router.HandlerError {
	managedRestaurants, handlerErr := getRestaurantsManagedThroughFB(userAccessToken, restaurants, fbAuth)
	if handlerErr != nil {
		return handlerErr
	}
	pageAccessTokens, handlerErr := getPageAccessTokensForRestaurants(userAccessToken, managedRestaurants, fbAuth)
	if handlerErr != nil {
		return handlerErr
	}
	if err := users.SetPageAccessTokens(fbUserID, pageAccessTokens); err != nil {
		return router.NewHandlerError(err, "Failed to persist Facebook page access tokens", http.StatusInternalServerError)
	}
	return nil
}

func getPageAccessTokensForRestaurants(userAccessToken *oauth2.Token, restaurants []*model.Restaurant, fbAuth facebook.Authenticator) ([]model.FacebookPageToken,
	*router.HandlerError) {
	pageAccessTokens := make([]model.FacebookPageToken, len(restaurants))
	for i, restaurant := range restaurants {
		pageAccessToken, err := fbAuth.PageAccessToken(userAccessToken, restaurant.FacebookPageID)
		if err == facebook.ErrNoSuchPage {
			return nil, router.NewHandlerError(err, "Access denied by Facebook to the managed page", http.StatusForbidden)
		} else if err != nil {
			return nil, router.NewHandlerError(err, "Failed to get access to the Facebook page", http.StatusInternalServerError)
		}
		pageAccessTokens[i] = model.FacebookPageToken{
			PageID: restaurant.FacebookPageID,
			Token:  pageAccessToken,
		}
	}
	return pageAccessTokens, nil
}

func getRestaurantsManagedThroughFB(userAccessToken *oauth2.Token, restaurants db.Restaurants, fbAuth facebook.Authenticator) ([]*model.Restaurant, *router.HandlerError) {
	fbPagesManagedByUser, handlerErr := getPages(userAccessToken, fbAuth)
	if handlerErr != nil {
		return nil, handlerErr
	}
	fbPageIDs := make([]string, len(fbPagesManagedByUser))
	for i, fbPage := range fbPagesManagedByUser {
		fbPageIDs[i] = fbPage.ID
	}
	restaurantsManagedByUserThroughFB, err := restaurants.GetByFacebookPageIDs(fbPageIDs)
	if err != nil {
		return nil, router.NewHandlerError(err, "Failed to find restaurants for FB pages associated with this user", http.StatusInternalServerError)
	}
	return restaurantsManagedByUserThroughFB, nil
}

func getLongTermToken(session string, r *http.Request, auther facebook.Authenticator) (*oauth2.Token, *router.HandlerError) {
	tok, err := auther.Token(session, r)
	if err != nil {
		if err == facebook.ErrMissingState {
			return nil, router.NewHandlerError(err, "Expecting a 'state' value", http.StatusBadRequest)
		} else if err == facebook.ErrInvalidState {
			return nil, router.NewHandlerError(err, "Invalid 'state' value", http.StatusForbidden)
		} else if err == facebook.ErrMissingCode {
			return nil, router.NewHandlerError(err, "Expecting a 'code' value", http.StatusBadRequest)
		}
		return nil, router.NewHandlerError(err, "Failed to connect to Facebook", http.StatusInternalServerError)
	}
	return tok, nil
}

func storeAccessTokensInDB(userID bson.ObjectId, fbUserID string, tok *oauth2.Token, sessionID string, usersCollection db.Users) *router.HandlerError {
	if err := usersCollection.SetAccessToken(fbUserID, *tok); err != nil {
		return router.NewHandlerError(err, "Failed to persist Facebook user access token in DB", http.StatusInternalServerError)
	}
	if err := usersCollection.SetSessionID(userID, sessionID); err != nil {
		return router.NewHandlerError(err, "Failed to persist session ID in DB", http.StatusInternalServerError)
	}
	return nil
}

func getUserID(tok *oauth2.Token, auther facebook.Authenticator) (string, error) {
	api := auther.APIConnection(tok)
	user, err := api.Me()
	if err != nil {
		return "", err
	}
	return user.ID, nil
}
