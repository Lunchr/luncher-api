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
func RedirectedFromFBForLogin(sessionManager session.Manager, auther facebook.Authenticator, usersCollection db.Users, restaurantsCollection db.Restaurants) router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		session := sessionManager.GetOrInit(w, r)
		tok, handlerErr := getLongTermToken(session, r, auther)
		if handlerErr != nil {
			return handlerErr
		}
		fbUserID, err := getUserID(tok, auther)
		if err != nil {
			return router.NewHandlerError(err, "Failed to get the user information from Facebook", http.StatusInternalServerError)
		}
		user, err := usersCollection.GetFbID(fbUserID)
		if err == mgo.ErrNotFound {
			return router.NewHandlerError(err, "User not registered", http.StatusForbidden)
		} else if err != nil {
			return router.NewHandlerError(err, "Failed to find the user from DB", http.StatusInternalServerError)
		}
		err = storeAccessTokensInDB(user.ID, fbUserID, tok, session, usersCollection)
		if err != nil {
			return router.NewHandlerError(err, "Failed to persist Facebook login information", http.StatusInternalServerError)
		}
		pageID, handlerErr := getPageID(user, restaurantsCollection)
		if handlerErr != nil {
			return handlerErr
		}
		if pageID != "" {
			pageAccessToken, err := auther.PageAccessToken(tok, pageID)
			if err != nil {
				if err == facebook.ErrNoSuchPage {
					return router.NewHandlerError(err, "Access denied by Facebook to the managed page", http.StatusForbidden)
				}
				return router.NewHandlerError(err, "Failed to get access to the Facebook page", http.StatusInternalServerError)
			}
			err = usersCollection.SetPageAccessToken(fbUserID, pageAccessToken)
			if err != nil {
				return router.NewHandlerError(err, "Failed to persist Facebook login information", http.StatusInternalServerError)
			}
		}
		http.Redirect(w, r, "/#/admin", http.StatusSeeOther)
		return nil
	}
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

func storeAccessTokensInDB(userID bson.ObjectId, fbUserID string, tok *oauth2.Token, sessionID string, usersCollection db.Users) error {
	err := usersCollection.SetAccessToken(fbUserID, *tok)
	if err != nil {
		return err
	}
	return usersCollection.SetSessionID(userID, sessionID)
}

func getUserID(tok *oauth2.Token, auther facebook.Authenticator) (string, error) {
	api := auther.APIConnection(tok)
	user, err := api.Me()
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

func getPageID(user *model.User, restaurantsCollection db.Restaurants) (string, *router.HandlerError) {
	restaurant, err := restaurantsCollection.GetID(user.RestaurantIDs[0])
	if err != nil {
		return "", router.NewHandlerError(err, "Failed to find the restaurant in the DB", http.StatusInternalServerError)
	}
	return restaurant.FacebookPageID, nil
}
