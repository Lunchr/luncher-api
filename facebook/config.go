package facebook

import "github.com/deiwin/gonfigure"

var (
	appIDEnvProperty     = gonfigure.NewRequiredEnvProperty("FACEBOOK_APP_ID")
	appSecretEnvProperty = gonfigure.NewRequiredEnvProperty("FACEBOOK_APP_SECRET")
)

type Config struct {
	AppID       string
	AppSecret   string
	RedirectURL string
	Scopes      []string
}

func NewConfig(redirectURL string, scopes []string) Config {
	return Config{
		AppID:       appIDEnvProperty.Value(),
		AppSecret:   appSecretEnvProperty.Value(),
		RedirectURL: redirectURL,
		Scopes:      scopes,
	}
}
