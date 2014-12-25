package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/deiwin/luncher-api/facebook"
	"github.com/deiwin/luncher-api/facebook/model"
	"github.com/deiwin/luncher-api/session"
)

type Facebook interface {
	// Login returns a handler that redirects the user to Facebook to log in
	Login() Handler
	// Redirected returns a handler that receives the user and page tokens for the
	// user who has just logged in through Facebook
	Redirected() Handler
}

type fbook struct {
	auth           facebook.Authenticator
	sessionManager session.Manager
	api            facebook.API
}

func NewFacebook(fbAuth facebook.Authenticator, sessMgr session.Manager, api facebook.API) Facebook {
	return fbook{fbAuth, sessMgr, api}
}

func (fb fbook) Login() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		session := fb.sessionManager.GetOrInitSession(w, r)
		redirectURL := fb.auth.AuthURL(session)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func (fb fbook) Redirected() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fb.checkState(w, r)
		if err != nil {
			log.Print(err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		code := r.FormValue("code")
		if code == "" {
			log.Println("A Facebook redirect request is missing the 'code' value")
			http.Error(w, "Expecting a 'code' value", http.StatusBadRequest)
			return
		}
		tok, err := fb.auth.Token(code)
		if err != nil {
			log.Print(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		client := fb.auth.Client(tok)

		connection := facebook.NewConnection(fb.api, client)
		accs, err := connection.Accounts()
		if err != nil {
			log.Print(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		pageID := "" // TODO get from DB. This should have been recorded on account creation/linking
		pageAccessToken, err := getPageAccessToken(accs, pageID)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, pageAccessToken)
	}
}

func getPageAccessToken(accs model.Accounts, pageID string) (pageAccessToken string, err error) {
	for _, page := range accs.Data {
		if page.ID == pageID {
			pageAccessToken = page.AccessToken
			return
		}
	}
	err = errors.New("Couldn't find the administered page")
	return
}

func (fb fbook) checkState(w http.ResponseWriter, r *http.Request) error {
	session := fb.sessionManager.GetOrInitSession(w, r)
	state := r.FormValue("state")
	if state == "" {
		return errors.New("A Facebook redirect request is missing the 'state' value")
	} else if state != session {
		return errors.New("A Facebook redirect request's 'state' value does not match the session")
	}
	return nil
}
