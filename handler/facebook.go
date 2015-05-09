package handler

import (
	"github.com/deiwin/facebook"
	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/session"
)

// Facebook holds the the handlers dealing with login and registration via Facebook
type Facebook struct {
	loginAuth        facebook.Authenticator
	registrationAuth facebook.Authenticator
	sessionManager   session.Manager
	usersCollection  db.Users
}

// NewFacebook initializes a Facebook struct
func NewFacebook(loginAuth facebook.Authenticator, registrationAuth facebook.Authenticator, sessMgr session.Manager, usersCollection db.Users) Facebook {
	return Facebook{loginAuth, registrationAuth, sessMgr, usersCollection}
}
