package handler

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	"gopkg.in/mgo.v2"

	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/router"
	fbmodel "github.com/deiwin/facebook/model"
	"github.com/julienschmidt/httprouter"
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

// RedirectedFromFBForRegistration provides a handler that stores the data about the current user
// required to continue the registration in the DB.
func (f Facebook) RedirectedFromFBForRegistration() router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		session := f.sessionManager.GetOrInit(w, r)
		tok, handlerErr := f.getLongTermToken(session, r)
		if handlerErr != nil {
			return handlerErr
		}
		fbUserID, err := f.getUserID(tok)
		if err != nil {
			return router.NewHandlerError(err, "Failed to get the user information from Facebook", http.StatusInternalServerError)
		}
		// We can't guarantee that the user doesn't just close the browser or something during the registration process.
		// Because of this, there already might be a user object with this FB User ID in the DB. If the user exists in the DB,
		// but doesn't have a restaurant connected to it, continue as the user would've been just created. But, to be safe, when
		// the user already exists and has a restaurant connected to it, fail immediately.
		user, err := f.usersCollection.GetFbID(fbUserID)
		if err == mgo.ErrNotFound {
			err = f.usersCollection.Insert(&model.User{FacebookUserID: fbUserID})
			if err != nil {
				return router.NewHandlerError(err, "Failed to create a User object in the DB", http.StatusInternalServerError)
			}
		} else if err != nil {
			return router.NewHandlerError(err, "Failed to check the DB for users", http.StatusInternalServerError)
		} else if err == nil && user.RestaurantID != "" {
			return router.NewSimpleHandlerError("This Facebook user is already registered", http.StatusForbidden)
		}
		err = f.storeAccessTokensInDB(fbUserID, tok, session)
		if err != nil {
			return router.NewHandlerError(err, "Failed to persist Facebook login information", http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/#/register/pages", http.StatusSeeOther)
		return nil
	}
}

// ListPagesManagedByUser returns a handler that lists all pages managed by the currently logged in user
func (f Facebook) ListPagesManagedByUser() router.Handler {
	handler := func(w http.ResponseWriter, r *http.Request, user *model.User) *router.HandlerError {
		fbPages, err := f.getPages(&user.Session.FacebookUserToken)
		if err != nil {
			return router.NewHandlerError(err, "Couldn't get the list of pages managed by this user", http.StatusBadGateway)
		}
		pages := mapFBPagesToModelPages(fbPages)
		return writeJSON(w, pages)
	}
	return checkLogin(f.sessionManager, f.usersCollection, handler)
}

// Page returns a handler that returns information about the specified page
func (f Facebook) Page() router.HandlerWithParams {
	handler := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, user *model.User) *router.HandlerError {
		id := ps.ByName("id")
		if id == "" {
			return router.NewSimpleHandlerError("Expecting a Page ID", http.StatusBadRequest)
		}
		fbPage, err := f.getPage(id, &user.Session.FacebookUserToken)
		if err != nil {
			return err
		}
		page := mapFBPageToModelPage(fbPage)
		return writeJSON(w, page)
	}
	return checkLoginWithParams(f.sessionManager, f.usersCollection, handler)
}

func (f Facebook) getPages(tok *oauth2.Token) ([]fbmodel.Page, error) {
	api := f.registrationAuth.APIConnection(tok)
	accs, err := api.Accounts()
	if err != nil {
		return nil, err
	}
	return accs.Data, nil
}

func (f Facebook) getPage(id string, tok *oauth2.Token) (*fbmodel.Page, *router.HandlerError) {
	api := f.registrationAuth.APIConnection(tok)
	page, err := api.Page(id)
	if err != nil {
		return nil, router.NewHandlerError(err, "Couldn't get the page", http.StatusBadGateway)
	}
	if page == nil {
		return nil, router.NewSimpleHandlerError("Page not found", http.StatusNotFound)
	}
	return page, nil
}

// Page represents a Facebook page
type Page struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Website string `json:"website,omitempty"`
}

// mapFBPagesToModelPages only maps the ID and the Name field, because that's really
// all that's required for the page listing
func mapFBPagesToModelPages(fbPages []fbmodel.Page) []Page {
	pages := make([]Page, len(fbPages))
	for i, v := range fbPages {
		pages[i] = Page{
			ID:   v.ID,
			Name: v.Name,
		}
	}
	return pages
}

func mapFBPageToModelPage(fbPage *fbmodel.Page) *Page {
	return &Page{
		ID:      fbPage.ID,
		Name:    fbPage.Name,
		Address: formAddressFromFBLocation(fbPage.Location),
		Phone:   fbPage.Phone,
		Website: fbPage.Website,
	}
}

func formAddressFromFBLocation(loc fbmodel.Location) string {
	return fmt.Sprintf("%s, %s, %s", loc.Street, loc.City, loc.Country)
}
