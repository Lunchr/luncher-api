package facebook

import "github.com/deiwin/gonfigure"

var (
	appIDEnvProperty     = gonfigure.NewRequiredEnvProperty("LUNCHER_FACEBOOK_APP_ID")
	appSecretEnvProperty = gonfigure.NewRequiredEnvProperty("LUNCHER_FACEBOOK_APP_SECRET")
)

type Config struct {
	AppID      string
	AppSecret  string
	ApiVersion string
}

func NewConfig() Config {
	return Config{
		AppID:       appIDEnvProperty.Value(),
		AppSecret:   appSecretEnvProperty.Value(),
		ApiVersion: apiVersionEnvProperty.DefaultValue(),
	}
}
