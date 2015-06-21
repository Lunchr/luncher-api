package handler

import (
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/router"
)

// RedirectToFBForRegistration returns a handler that redirects the user to Facebook to log in
// so they could be registered in our system
func (f Facebook) RedirectToFBForRegistration() router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		session := f.sessionManager.GetOrInit(w, r)
		redirectURL := f.registrationAuth.AuthURL(session)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return nil
	}
}

// RedirectedFromFBForRegistration provides a handler that stores the data about the current user
// required to continue the registration in the DB.
func (f Facebook) RedirectedFromFBForRegistration() router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		session := f.sessionManager.GetOrInit(w, r)
		tok, handlerErr := f.getLongTermToken(session, r)
		if handlerErr != nil {
			return handlerErr
		}
		fbUserID, err := f.getUserID(tok)
		if err != nil {
			return router.NewHandlerError(err, "Failed to get the user information from Facebook", http.StatusInternalServerError)
		}
		// We can't guarantee that the user doesn't just close the browser or something during the registration process.
		// Because of this, there already might be a user object with this FB User ID in the DB.
		_, err = f.usersCollection.GetFbID(fbUserID)
		if err == mgo.ErrNotFound {
			err = f.usersCollection.Insert(&model.User{FacebookUserID: fbUserID})
			if err != nil {
				return router.NewHandlerError(err, "Failed to create a User object in the DB", http.StatusInternalServerError)
			}
		} else if err != nil {
			return router.NewHandlerError(err, "Failed to check the DB for users", http.StatusInternalServerError)
		}
		err = f.storeAccessTokensInDB(fbUserID, tok, session)
		if err != nil {
			return router.NewHandlerError(err, "Failed to persist Facebook login information", http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/#/register/pages", http.StatusSeeOther)
		return nil
	}
}
