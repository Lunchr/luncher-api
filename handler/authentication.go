package handler

import (
	"net/http"

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
		user, err := getUserForSession(sessionManager, usersCollection, r)
		if err != nil {
			return router.NewHandlerError(err, "", http.StatusForbidden)
		}
		return handler(w, r, user)
	}
}

func checkLoginWithParams(sessionManager session.Manager, usersCollection db.Users, handler HandlerWithParamsWithUser) router.HandlerWithParams {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) *router.HandlerError {
		user, err := getUserForSession(sessionManager, usersCollection, r)
		if err != nil {
			return router.NewHandlerError(err, "", http.StatusForbidden)
		}
		return handler(w, r, ps, user)
	}
}

func getUserForSession(sessionManager session.Manager, usersCollection db.Users, r *http.Request) (*model.User, error) {
	session, err := sessionManager.Get(r)
	if err != nil {
		return nil, err
	}
	user, err := usersCollection.GetSessionID(session)
	if err != nil {
		return nil, err
	}
	return user, nil
}
