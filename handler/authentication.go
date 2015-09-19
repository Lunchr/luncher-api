package handler

import (
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"github.com/julienschmidt/httprouter"
)

func Logout(sessionManager session.Manager, usersCollection db.Users) router.Handler {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User) *router.HandlerError {
		if err := usersCollection.UnsetSessionID(user.ID); err != nil {
			return router.NewHandlerError(err, "", http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	}
	return checkLogin(sessionManager, usersCollection, handler)
}

type HandlerWithUser func(w http.ResponseWriter, r *http.Request, user *model.User) *router.HandlerError
type HandlerWithParamsWithUser func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user *model.User) *router.HandlerError

func checkLogin(sessionManager session.Manager, usersCollection db.Users, handler HandlerWithUser) router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		user, handlerErr := getUserForSession(sessionManager, usersCollection, r)
		if handlerErr != nil {
			return handlerErr
		}
		return handler(w, r, user)
	}
}

func checkLoginWithParams(sessionManager session.Manager, usersCollection db.Users, handler HandlerWithParamsWithUser) router.HandlerWithParams {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) *router.HandlerError {
		user, handlerErr := getUserForSession(sessionManager, usersCollection, r)
		if handlerErr != nil {
			return handlerErr
		}
		return handler(w, r, ps, user)
	}
}

func getUserForSession(sessionManager session.Manager, usersCollection db.Users, r *http.Request) (*model.User, *router.HandlerError) {
	sessionID, err := sessionManager.Get(r)
	if err == session.ErrNotFound {
		return nil, router.NewHandlerError(err, "Session not established for this connection", http.StatusUnauthorized)
	} else if err != nil {
		return nil, router.NewHandlerError(err, "Failed to get the session", http.StatusInternalServerError)
	}
	user, err := usersCollection.GetSessionID(sessionID)
	if err == mgo.ErrNotFound {
		return nil, router.NewHandlerError(err, "User not logged in", http.StatusUnauthorized)
	} else if err != nil {
		return nil, router.NewHandlerError(err, "Failed to find the user for this session", http.StatusInternalServerError)
	}
	return user, nil
}
