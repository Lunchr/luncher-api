package db

import "github.com/deiwin/gonfigure"

var (
	dbURLEnvProperty  = gonfigure.NewEnvProperty("LUNCHER_DB_ADDRESS", "mongodb://localhost/test")
	dbNameEnvProperty = gonfigure.NewEnvProperty("LUNCHER_DB_NAME", "")
)

type Config struct {
	DbURL  string
	DbName string
}

func NewConfig() *Config {
	return &Config{
		DbURL:  dbURLEnvProperty.Value(),
		DbName: dbNameEnvProperty.Value(),
	}
}
