package facebook

import "github.com/deiwin/luncher-api/config"

var (
	appIDEnvProperty     = config.NewEnvProperty("LUNCHER_FACEBOOK_APP_ID", "1")
	appSecretEnvProperty = config.NewEnvProperty("LUNCHER_FACEBOOK_APP_SECRET", "secret")
)

type Config struct {
	AppID     string
	AppSecret string
}

func NewConfig() *Config {
	return &Config{
		AppID:     appIDEnvProperty.DefaultValue(),
		AppSecret: appSecretEnvProperty.DefaultValue(),
	}
}
