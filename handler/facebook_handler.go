package handler

import (
	"net/http"

	"github.com/deiwin/facebook"
	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"golang.org/x/oauth2"
)

type Facebook interface {
	// Login returns a handler that redirects the user to Facebook to log in
	Login() router.Handler
	// Redirected returns a handler that receives the user and page tokens for the
	// user who has just logged in through Facebook. Updates the user and page
	// access tokens in the DB
	Redirected() router.Handler
}

type fbook struct {
	auth            facebook.Authenticator
	sessionManager  session.Manager
	usersCollection db.Users
}

func NewFacebook(fbAuth facebook.Authenticator, sessMgr session.Manager, usersCollection db.Users) Facebook {
	return fbook{fbAuth, sessMgr, usersCollection}
}

func (fb fbook) Login() router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		session := fb.sessionManager.GetOrInit(w, r)
		redirectURL := fb.auth.AuthURL(session)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return nil
	}
}

func (fb fbook) Redirected() router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		session := fb.sessionManager.GetOrInit(w, r)
		tok, err := fb.auth.Token(session, r)
		if err != nil {
			if err == facebook.ErrMissingState {
				return router.NewHandlerError(err, "Expecting a 'state' value", http.StatusBadRequest)
			} else if err == facebook.ErrInvalidState {
				return router.NewHandlerError(err, "Invalid 'state' value", http.StatusForbidden)
			} else if err == facebook.ErrMissingCode {
				return router.NewHandlerError(err, "Expecting a 'code' value", http.StatusBadRequest)
			}
			return router.NewHandlerError(err, "Failed to connect to Facebook", http.StatusInternalServerError)
		}
		fbUserID, err := fb.getUserID(tok)
		if err != nil {
			return router.NewHandlerError(err, "Failed to get the user information from Facebook", http.StatusInternalServerError)
		}
		err = fb.storeAccessTokensInDB(fbUserID, tok, session)
		if err != nil {
			return router.NewHandlerError(err, "Failed to persist Facebook login information", http.StatusInternalServerError)
		}
		pageID, handlerErr := fb.getPageID(fbUserID)
		if handlerErr != nil {
			return handlerErr
		}
		if pageID != "" {
			pageAccessToken, err := fb.auth.PageAccessToken(tok, pageID)
			if err != nil {
				if err == facebook.ErrNoSuchPage {
					return router.NewHandlerError(err, "Access denied by Facebook to the managed page", http.StatusForbidden)
				}
				return router.NewHandlerError(err, "Failed to get access to the Facebook page", http.StatusInternalServerError)
			}
			err = fb.usersCollection.SetPageAccessToken(fbUserID, pageAccessToken)
			if err != nil {
				return router.NewHandlerError(err, "Failed to persist Facebook login information", http.StatusInternalServerError)
			}
		}
		http.Redirect(w, r, "/#/admin", http.StatusSeeOther)
		return nil
	}
}

func (fb fbook) storeAccessTokensInDB(fbUserID string, tok *oauth2.Token, sessionID string) error {
	err := fb.usersCollection.SetAccessToken(fbUserID, *tok)
	if err != nil {
		return err
	}
	user, err := fb.usersCollection.GetFbID(fbUserID)
	if err != nil {
		return err
	}
	return fb.usersCollection.SetSessionID(user.ID, sessionID)
}

func (fb fbook) getUserID(tok *oauth2.Token) (string, error) {
	api := fb.auth.APIConnection(tok)
	user, err := api.Me()
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

func (fb fbook) getPageID(userID string) (string, *router.HandlerError) {
	userInDB, err := fb.usersCollection.GetFbID(userID)
	if err != nil {
		return "", router.NewHandlerError(err, "Failed to find a user in DB related to this Facebook User ID", http.StatusInternalServerError)
	}
	return userInDB.FacebookPageID, nil
}
