package geo

import "github.com/deiwin/gonfigure"

var (
	apiKeyProperty = gonfigure.NewRequiredEnvProperty("GOOGLE_GEOCODING_API_KEY")
)

type Config struct {
	APIKey string
}

func NewConfig() *Config {
	return &Config{
		APIKey: apiKeyProperty.Value(),
	}
}
