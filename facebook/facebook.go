package facebook

import (
	"log"

	"golang.org/x/oauth2"
)

// Authenticator provides the authentication functionality for Facebook users
// using Facebook's OAuth
type Authenticator interface {
	// AuthURL returns a Facebook URL the user should be redirect to. The user
	// will then be asked to log in by Facebook at that URL and will be redirected
	// back to our API by Facebook.
	AuthURL(session string) string
}

func NewAuthenticator(conf Config, domain string) Authenticator {
	return authenticator{conf, domain}
}

type authenticator struct {
	conf   Config
	domain string
}

func (a authenticator) AuthURL(session string) string {
	opts, err := oauth2.New(
		oauth2.Client(a.conf.AppID, a.conf.AppSecret),
		oauth2.RedirectURL(a.domain+"api/v1/oauth/facebook/redirect"),
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
