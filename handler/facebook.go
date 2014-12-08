package handler

import (
	"net/http"

	"github.com/deiwin/luncher-api/facebook"
	"github.com/deiwin/luncher-api/session"
)

// FacebookLogin returns a handler that redirects the user to Facebook to log in
func FacebookLogin(fbAuth facebook.Authenticator, sessMgr session.Manager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := sessMgr.GetOrInitSession(w, r)
		redirectURL := fbAuth.AuthURL(session)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}
