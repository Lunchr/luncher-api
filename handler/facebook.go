package handler

import (
	"log"
	"net/http"

	"github.com/deiwin/luncher-api/facebook"
	"github.com/deiwin/luncher-api/session"
)

type Facebook interface {
	// Login returns a handler that redirects the user to Facebook to log in
	Login() handler
	// Redirected returns a handler that TODO
	Redirected() handler
}

type fbook struct {
	auth           facebook.Authenticator
	sessionManager session.Manager
	api            facebook.API
}

func NewFacebook(fbAuth facebook.Authenticator, sessMgr session.Manager, api facebook.API) Facebook {
	return fbook{fbAuth, sessMgr, api}
}

func (fb fbook) Login() handler {
	return func(w http.ResponseWriter, r *http.Request) {
		session := fb.sessionManager.GetOrInitSession(w, r)
		redirectURL := fb.auth.AuthURL(session)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func (fb fbook) Redirected() handler {
	return func(w http.ResponseWriter, r *http.Request) {
		session := fb.sessionManager.GetOrInitSession(w, r)
		state := r.FormValue("state")
		if state == "" {
			log.Println("A Facebook redirect request is missing the 'state' value")
			http.Error(w, "Expecting a 'state' value", http.StatusBadRequest)
			return
		} else if state != session {
			log.Println(state)
			log.Println(session)
			log.Println("A Facebook redirect request's 'state' value does not match the session")
			http.Error(w, "Wrong 'state' value", http.StatusBadRequest)
			return
		}

		code := r.FormValue("code")
		if code == "" {
			log.Println("A Facebook redirect request is missing the 'code' value")
			http.Error(w, "Expecting a 'code' value", http.StatusBadRequest)
			return
		}
		client, err := fb.auth.CreateClient(code)
		if err != nil {
			log.Print(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		connection := facebook.NewConnection(fb.api, client)
		user, err := connection.Me()
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("id:" + user.Id))
	}
}
