package handler

import (
	"net/http"

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

// TODO
func (f Facebook) RedirectedFromFBForRegistration() router.Handler {
	return nil
}
