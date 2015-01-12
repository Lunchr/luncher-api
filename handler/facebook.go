package handler

import (
	"errors"
	"net/http"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/facebook"
	"github.com/deiwin/luncher-api/session"
	"golang.org/x/oauth2"
)

type Facebook interface {
	// Login returns a handler that redirects the user to Facebook to log in
	Login() Handler
	// Redirected returns a handler that receives the user and page tokens for the
	// user who has just logged in through Facebook. Updates the user and page
	// access tokens in the DB
	Redirected() Handler
}

type fbook struct {
	auth            facebook.Authenticator
	sessionManager  session.Manager
	api             facebook.API
	usersCollection db.Users
}

func NewFacebook(fbAuth facebook.Authenticator, sessMgr session.Manager, api facebook.API, usersCollection db.Users) Facebook {
	return fbook{fbAuth, sessMgr, api, usersCollection}
}

func (fb fbook) Login() Handler {
	return func(w http.ResponseWriter, r *http.Request) *handlerError {
		session := fb.sessionManager.GetOrInitSession(w, r)
		redirectURL := fb.auth.AuthURL(session)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return nil
	}
}

func (fb fbook) Redirected() Handler {
	return func(w http.ResponseWriter, r *http.Request) *handlerError {
		if err := fb.checkState(w, r); err != nil {
			return err
		}
		code := r.FormValue("code")
		if code == "" {
			err := errors.New("A Facebook redirect request is missing the 'code' value")
			return &handlerError{err, "Expecting a 'code' value", http.StatusBadRequest}
		}
		tok, err := fb.auth.Token(code)
		if err != nil {
			return &handlerError{err, "", http.StatusInternalServerError}
		}
		client := fb.auth.Client(tok)
		connection := facebook.NewConnection(fb.api, client)
		userID, err := getUserID(connection)
		if err != nil {
			return &handlerError{err, "", http.StatusInternalServerError}
		}
		pageID, err := fb.getPageID(userID)
		if err != nil {
			return &handlerError{err, "", http.StatusInternalServerError}
		}
		pageAccessToken, err := fb.getPageAccessToken(connection, pageID)
		if err != nil {
			return &handlerError{err, "", http.StatusInternalServerError}
		}
		err = fb.storeAccessTokensInDB(userID, tok, pageAccessToken)
		if err != nil {
			return &handlerError{err, "", http.StatusInternalServerError}
		}
		http.Redirect(w, r, "/#/admin", http.StatusSeeOther)
		return nil
	}
}

func (fb fbook) storeAccessTokensInDB(userID string, tok *oauth2.Token, pageAccessToken string) (err error) {
	err = fb.usersCollection.SetAccessToken(userID, *tok)
	if err != nil {
		return
	}
	err = fb.usersCollection.SetPageAccessToken(userID, pageAccessToken)
	return
}

func (fb fbook) getPageAccessToken(connection facebook.Connection, pageID string) (pageAccessToken string, err error) {
	accs, err := connection.Accounts()
	if err != nil {
		return
	}
	for _, page := range accs.Data {
		if page.ID == pageID {
			pageAccessToken = page.AccessToken
			return
		}
	}
	err = errors.New("Couldn't find the administered page")
	return
}

func (fb fbook) checkState(w http.ResponseWriter, r *http.Request) *handlerError {
	session := fb.sessionManager.GetOrInitSession(w, r)
	state := r.FormValue("state")
	if state == "" {
		err := errors.New("A Facebook redirect request is missing the 'state' value")
		return &handlerError{err, "Expecting a 'state' value", http.StatusBadRequest}
	} else if state != session {
		err := errors.New("A Facebook redirect request's 'state' value does not match the session")
		return &handlerError{err, "Invalid 'state' value", http.StatusForbidden}
	}
	return nil
}

func getUserID(connection facebook.Connection) (userID string, err error) {
	user, err := connection.Me()
	if err != nil {
		return
	}
	userID = user.Id
	return
}

func (fb fbook) getPageID(userID string) (pageID string, err error) {
	userInDB, err := fb.usersCollection.Get(userID)
	if err != nil {
		return
	}
	pageID = userInDB.FacebookPageID
	return
}
