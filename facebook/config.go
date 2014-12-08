package facebook

import "github.com/deiwin/luncher-api/config"

var (
	appIDEnvProperty     = config.NewRequiredEnvProperty("LUNCHER_FACEBOOK_APP_ID")
	appSecretEnvProperty = config.NewRequiredEnvProperty("LUNCHER_FACEBOOK_APP_SECRET")
)

type Config struct {
	AppID     string
	AppSecret string
}

func NewConfig() Config {
	return Config{
		AppID:     appIDEnvProperty.DefaultValue(),
		AppSecret: appSecretEnvProperty.DefaultValue(),
	}
}
