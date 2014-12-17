package handler

import (
	"log"
	"net/http"

	"github.com/deiwin/luncher-api/facebook"
	"github.com/deiwin/luncher-api/session"
)

// FacebookLogin returns a handler that redirects the user to Facebook to log in
func FacebookLogin(fbAuth facebook.Authenticator, sessMgr session.Manager) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		session := sessMgr.GetOrInitSession(w, r)
		redirectURL := fbAuth.AuthURL(session)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func FacebookRedirected(fbAuth facebook.Authenticator, sessMgr session.Manager) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		session := sessMgr.GetOrInitSession(w, r)
		state := r.FormValue("state")
		if state == "" {
			log.Println("A Facebook redirect request is missing the 'state' value")
			http.Error(w, "Expecting a 'state' value", http.StatusBadRequest)
			return
		} else if state != session {
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
		transport, err := fbAuth.CreateTransport(code)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
