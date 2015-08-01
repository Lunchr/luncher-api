package handler

import (
	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/session"
	"github.com/deiwin/facebook"
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
