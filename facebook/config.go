package facebook

import "github.com/deiwin/luncher-api/config"

var (
	appIDEnvProperty      = config.NewRequiredEnvProperty("LUNCHER_FACEBOOK_APP_ID")
	appSecretEnvProperty  = config.NewRequiredEnvProperty("LUNCHER_FACEBOOK_APP_SECRET")
	apiVersionEnvProperty = config.NewEnvProperty("LUNCHER_FACEBOOK_API_VERSION", "v2.2")
)

type Config struct {
	AppID      string
	AppSecret  string
	ApiVersion string
}

func NewConfig() Config {
	return Config{
		AppID:      appIDEnvProperty.DefaultValue(),
		AppSecret:  appSecretEnvProperty.DefaultValue(),
		ApiVersion: apiVersionEnvProperty.DefaultValue(),
	}
}
