package facebook

import (
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
	Token(code string) (*oauth2.Token, error)
	// Client returns an *http.Client that can be used to make authenticated
	// requests to the Facebook API.
	Client(tok *oauth2.Token) *http.Client
}

// NewAuthenticator initializes and returns an Authenticator
func NewAuthenticator(conf Config, domain string) Authenticator {
	opts := &oauth2.Config{
		ClientID:     conf.AppID,
		ClientSecret: conf.AppSecret,
		RedirectURL:  domain + "/api/v1/login/facebook/redirected",
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

func (a authenticator) Token(code string) (*oauth2.Token, error) {
	return a.Exchange(oauth2.NoContext, code)
}

func (a authenticator) Client(tok *oauth2.Token) *http.Client {
	return a.Config.Client(oauth2.NoContext, tok)
}
