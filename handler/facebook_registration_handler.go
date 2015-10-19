package handler

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	"gopkg.in/mgo.v2"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"github.com/deiwin/facebook"
	fbmodel "github.com/deiwin/facebook/model"
	"github.com/julienschmidt/httprouter"
)

// RedirectToFBForRegistration returns a handler that redirects the user to Facebook to log in
// so they could be registered in our system
func RedirectToFBForRegistration(sessionManager session.Manager, auther facebook.Authenticator,
	tokens db.RegistrationAccessTokens) router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		tokenString := r.FormValue("token")
		if tokenString != "" {
			token, err := model.TokenFromString(tokenString)
			if err != nil {
				return router.NewHandlerError(err, "Failed to parse the token", http.StatusBadRequest)
			}
			tokenExists, err := tokens.Exists(token)
			if err != nil {
				return router.NewHandlerError(err, "Failed to check if token exists in DB", http.StatusBadRequest)
			} else if !tokenExists {
				return router.NewSimpleHandlerError("Invalid access token", http.StatusForbidden)
			}
		}
		session := sessionManager.GetOrInit(w, r)
		redirectURL := auther.AuthURL(session)
		return writeString(w, redirectURL)
	}
}

// RedirectedFromFBForRegistration provides a handler that stores the data about the current user
// required to continue the registration in the DB.
func RedirectedFromFBForRegistration(sessionManager session.Manager, auther facebook.Authenticator, usersCollection db.Users) router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		session := sessionManager.GetOrInit(w, r)
		tok, handlerErr := getLongTermToken(session, r, auther)
		if handlerErr != nil {
			return handlerErr
		}
		fbUserID, err := getUserID(tok, auther)
		if err != nil {
			return router.NewHandlerError(err, "Failed to get the user information from Facebook", http.StatusBadGateway)
		}
		_, err = usersCollection.GetFbID(fbUserID)
		if err == nil {
			return router.NewSimpleHandlerError("This Facebook user is already registered", http.StatusBadRequest)
		} else if err != mgo.ErrNotFound {
			return router.NewHandlerError(err, "Failed to check the DB for users", http.StatusInternalServerError)
		}
		err = usersCollection.Insert(&model.User{FacebookUserID: fbUserID})
		if err != nil {
			return router.NewHandlerError(err, "Failed to create a User object in the DB", http.StatusInternalServerError)
		}
		// We're not storing the session data in the DB so that the user will be asked to log in again when they get
		// redirected to the admin page
		http.Redirect(w, r, "/#/admin", http.StatusSeeOther)
		return nil
	}
}

// ListPagesManagedByUser returns a handler that lists all pages managed by the currently logged in user
func ListPagesManagedByUser(sessionManager session.Manager, auther facebook.Authenticator, usersCollection db.Users) router.Handler {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User) *router.HandlerError {
		fbPages, handlerErr := getPages(&user.Session.FacebookUserToken, auther)
		if handlerErr != nil {
			return handlerErr
		}
		pages := mapFBPagesToModelPages(fbPages)
		var pagesNotAlreadyRegistered []FacebookPage
		for _, page := range pages {
			if !pageAlreadyRegistered(page.ID, user.Session.FacebookPageTokens) {
				pagesNotAlreadyRegistered = append(pagesNotAlreadyRegistered, page)
			}
		}
		return writeJSON(w, pagesNotAlreadyRegistered)
	}
	return checkLogin(sessionManager, usersCollection, handler)
}

func pageAlreadyRegistered(pageID string, accessTokensOfRegisteredPages []model.FacebookPageToken) bool {
	for _, pageAccessToken := range accessTokensOfRegisteredPages {
		if pageAccessToken.PageID == pageID {
			return true
		}
	}
	return false
}

// Page returns a handler that returns information about the specified page
func Page(sessionManager session.Manager, auther facebook.Authenticator, usersCollection db.Users) router.HandlerWithParams {
	handler := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user *model.User) *router.HandlerError {
		id := ps.ByName("id")
		if id == "" {
			return router.NewSimpleHandlerError("Expecting a Page ID", http.StatusBadRequest)
		}
		fbPage, err := getPage(id, &user.Session.FacebookUserToken, auther)
		if err != nil {
			return err
		}
		page := mapFBPageToModelPage(fbPage)
		return writeJSON(w, page)
	}
	return checkLoginWithParams(sessionManager, usersCollection, handler)
}

func getPages(tok *oauth2.Token, auther facebook.Authenticator) ([]fbmodel.Page, *router.HandlerError) {
	api := auther.APIConnection(tok)
	accs, err := api.Accounts()
	if err != nil {
		return nil, router.NewHandlerError(err, "Couldn't get the list of pages managed by this user", http.StatusBadGateway)
	}
	return accs.Data, nil
}

func getPage(id string, tok *oauth2.Token, auther facebook.Authenticator) (*fbmodel.Page, *router.HandlerError) {
	api := auther.APIConnection(tok)
	page, err := api.Page(id)
	if err != nil {
		return nil, router.NewHandlerError(err, "Couldn't get the page", http.StatusBadGateway)
	}
	if page == nil {
		return nil, router.NewSimpleHandlerError("Page not found", http.StatusNotFound)
	}
	return page, nil
}

// FacebookPage defines the response format for the Page() handler
type FacebookPage struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Website string `json:"website,omitempty"`
	Email   string `json:"emails,omitempty"`
}

// mapFBPagesToModelPages only maps the ID and the Name field, because that's really
// all that's required for the page listing
func mapFBPagesToModelPages(fbPages []fbmodel.Page) []FacebookPage {
	pages := make([]FacebookPage, len(fbPages))
	for i, v := range fbPages {
		pages[i] = FacebookPage{
			ID:   v.ID,
			Name: v.Name,
		}
	}
	return pages
}

func mapFBPageToModelPage(fbPage *fbmodel.Page) *FacebookPage {
	var page = FacebookPage{
		ID:      fbPage.ID,
		Name:    fbPage.Name,
		Address: formAddressFromFBLocation(fbPage.Location),
		Phone:   fbPage.Phone,
		Website: fbPage.Website,
	}
	if len(fbPage.Emails) >= 1 {
		// There could be multiple emails, but we'll just use the first one
		page.Email = fbPage.Emails[0]
	}
	return &page
}

func formAddressFromFBLocation(loc fbmodel.Location) string {
	return fmt.Sprintf("%s, %s, %s", loc.Street, loc.City, loc.Country)
}
