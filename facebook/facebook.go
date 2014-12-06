package facebook

import (
	"log"

	"golang.org/x/oauth2"
)

var domain = "haha" // XXX should this be in the FB conf?

func something(conf Config, session string) string {
	opts, err := oauth2.New(
		oauth2.Client(conf.AppID, conf.AppSecret),
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
