package db

import "github.com/deiwin/luncher-api/config"

var (
	dbURLEnvProperty  = config.NewEnvProperty("LUNCHER_DB_ADDRESS", "mongodb://localhost/test")
	dbNameEnvProperty = config.NewEnvProperty("LUNCHER_DB_NAME", "")
)

type Config struct {
	DbURL  string
	DbName string
}

func NewConfig() *Config {
	return &Config{
		DbURL:  dbURLEnvProperty.DefaultValue(),
		DbName: dbNameEnvProperty.DefaultValue(),
	}
}
