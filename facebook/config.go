package facebook

import "os"

const (
	appIDEnvVariable     = "LUNCHER_FACEBOOK_APP_ID"
	appSecretEnvVariable = "LUNCHER_FACEBOOK_APP_SECRET"
	defaultAppID         = "1"
	defaultAppSecret     = "secret"
)

type Config struct {
	AppID     string
	AppSecret string
}

func NewConfig() *Config {
	return &Config{
		AppID:     getEnvOrDefaultAppID(),
		AppSecret: getEnvOrDefaultAppSecret(),
	}
}

func getEnvOrDefaultAppID() (appID string) {
	appID = os.Getenv(appIDEnvVariable)
	if appID == "" {
		appID = defaultAppID
	}
	return
}

func getEnvOrDefaultAppSecret() (appSecret string) {
	appSecret = os.Getenv(appSecretEnvVariable)
	if appSecret == "" {
		appSecret = defaultAppSecret
	}
	return
}
