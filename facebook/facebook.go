package facebook

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

// Authenticator provides the authentication functionality for Facebook users
// using Facebook's OAuth
type Authenticator interface {
	// AuthURL returns a Facebook URL the user should be redirect to. The user
	// will then be asked to log in by Facebook at that URL and will be redirected
	// back to our API by Facebook.
	AuthURL(session string) string
	// CreateClient returns an *http.Client that can be used to make authenticated
	// requests to the Facebook API
	CreateClient(code string) (*http.Client, error)
}

// NewAuthenticator initializes and returns an Authenticator
func NewAuthenticator(conf Config, domain string) Authenticator {
	opts := &oauth2.Config{
		ClientID:     conf.AppID,
		ClientSecret: conf.AppSecret,
		RedirectURL:  domain + "api/v1/login/facebook/redirected",
		Scopes:       []string{"manage_pages", "publish_actions"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.facebook.com/dialog/oauth",
			TokenURL: "https://graph.facebook.com/oauth/access_token",
		},
	}
	return authenticator{opts}
}

type authenticator struct {
	*oauth2.Config
}

func (a authenticator) AuthURL(session string) string {
	return a.AuthCodeURL(session, oauth2.AccessTypeOffline)
}

func (a authenticator) CreateClient(code string) (client *http.Client, err error) {
	if _, err = fmt.Scan(&code); err != nil {
		return
	}
	tok, err := a.Exchange(oauth2.NoContext, code)
	if err != nil {
		return
	}
	client = a.Client(oauth2.NoContext, tok)
	return
}
