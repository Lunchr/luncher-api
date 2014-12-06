package facebook

import (
	"log"

	"golang.org/x/oauth2"
)

var domain = "haha" // XXX should this be in the FB conf?

// Authenticator provides the authentication functionality for Facebook users
// using Facebook's OAuth
type Authenticator interface {
	// AuthURL returns a Facebook URL the user should be redirect to. The user
	// will then be asked to log in by Facebook at that URL and will be redirected
	// back to our API by Facebook.
	AuthURL(session string) string
}

func NewAuthenticator(conf Config) Authenticator {
	return authenticator{conf}
}

type authenticator struct {
	conf Config
}

func (auth authenticator) AuthURL(session string) string {
	opts, err := oauth2.New(
		oauth2.Client(auth.conf.AppID, auth.conf.AppSecret),
		oauth2.RedirectURL(domain+"api/v1/oauth/facebook/redirect"),
		oauth2.Scope("manage_pages", "publish_actions"),
		oauth2.Endpoint(
			"https://www.facebook.com/dialog/oauth",
			"https://graph.facebook.com/oauth/access_token",
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	return opts.AuthCodeURL(session, "offline", "auto")
}
