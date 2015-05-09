package handler

import (
	"github.com/deiwin/facebook"
	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/router"
	"github.com/deiwin/luncher-api/session"
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
	loginAuth        facebook.Authenticator
	registrationAuth facebook.Authenticator
	sessionManager   session.Manager
	usersCollection  db.Users
}

func NewFacebook(loginAuth facebook.Authenticator, registrationAuth facebook.Authenticator, sessMgr session.Manager, usersCollection db.Users) Facebook {
	return fbook{loginAuth, registrationAuth, sessMgr, usersCollection}
}
