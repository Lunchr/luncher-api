package handler

import (
	"net/http"

	"github.com/deiwin/facebook"
	"github.com/deiwin/luncher-api/router"
	"golang.org/x/oauth2"
)

// RedirectToFBForLogin returns a handler that redirects the user to Facebook to log in
func (f Facebook) RedirectToFBForLogin() router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		session := f.sessionManager.GetOrInit(w, r)
		redirectURL := f.loginAuth.AuthURL(session)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return nil
	}
}

// RedirectedFromFBForLogin returns a handler that receives the user and page tokens for the
// user who has just logged in through Facebook. Updates the user and page
// access tokens in the DB
func (f Facebook) RedirectedFromFBForLogin() router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		session := f.sessionManager.GetOrInit(w, r)
		tok, err := f.loginAuth.Token(session, r)
		if err != nil {
			if err == facebook.ErrMissingState {
				return &router.HandlerError{err, "Expecting a 'state' value", http.StatusBadRequest}
			} else if err == facebook.ErrInvalidState {
				return &router.HandlerError{err, "Invalid 'state' value", http.StatusForbidden}
			} else if err == facebook.ErrMissingCode {
				return &router.HandlerError{err, "Expecting a 'code' value", http.StatusBadRequest}
			}
			return &router.HandlerError{err, "Failed to connect to Facebook", http.StatusInternalServerError}
		}
		fbUserID, err := f.getUserID(tok)
		if err != nil {
			return &router.HandlerError{err, "Failed to get the user information from Facebook", http.StatusInternalServerError}
		}
		err = f.storeAccessTokensInDB(fbUserID, tok, session)
		if err != nil {
			return &router.HandlerError{err, "Failed to persist Facebook login information", http.StatusInternalServerError}
		}
		pageID, handlerErr := f.getPageID(fbUserID)
		if handlerErr != nil {
			return handlerErr
		}
		if pageID != "" {
			pageAccessToken, err := f.loginAuth.PageAccessToken(tok, pageID)
			if err != nil {
				if err == facebook.ErrNoSuchPage {
					return &router.HandlerError{err, "Access denied by Facebook to the managed page", http.StatusForbidden}
				}
				return &router.HandlerError{err, "Failed to get access to the Facebook page", http.StatusInternalServerError}
			}
			err = f.usersCollection.SetPageAccessToken(fbUserID, pageAccessToken)
			if err != nil {
				return &router.HandlerError{err, "Failed to persist Facebook login information", http.StatusInternalServerError}
			}
		}
		http.Redirect(w, r, "/#/admin", http.StatusSeeOther)
		return nil
	}
}

func (f Facebook) storeAccessTokensInDB(fbUserID string, tok *oauth2.Token, sessionID string) error {
	err := f.usersCollection.SetAccessToken(fbUserID, *tok)
	if err != nil {
		return err
	}
	user, err := f.usersCollection.GetFbID(fbUserID)
	if err != nil {
		return err
	}
	return f.usersCollection.SetSessionID(user.ID, sessionID)
}

func (f Facebook) getUserID(tok *oauth2.Token) (string, error) {
	api := f.loginAuth.APIConnection(tok)
	user, err := api.Me()
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

func (f Facebook) getPageID(userID string) (string, *router.HandlerError) {
	userInDB, err := f.usersCollection.GetFbID(userID)
	if err != nil {
		return "", &router.HandlerError{err, "Failed to find a user in DB related to this Facebook User ID", http.StatusInternalServerError}
	}
	return userInDB.FacebookPageID, nil
}
