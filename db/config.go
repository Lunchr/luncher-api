package db

import "os"

const (
	dbURLEnvVariable  = "PRAAD_DB_ADDRESS"
	dbNameEnvVariable = "PRAAD_DB_NAME"
	defaultDbURL      = "mongodb://localhost/test"
	defaultDbName     = ""
)

type Config struct {
	DbURL  string
	DbName string
}

func NewConfig() *Config {
	return &Config{
		DbURL:  getEnvOrDefaultDbURL(),
		DbName: getEnvOrDefaultDbName(),
	}
}

func getEnvOrDefaultDbURL() (dbURL string) {
	dbURL = os.Getenv(dbURLEnvVariable)
	if dbURL == "" {
		dbURL = defaultDbURL
	}
	return
}

func getEnvOrDefaultDbName() (dbName string) {
	dbName = os.Getenv(dbNameEnvVariable)
	if dbName == "" {
		dbName = defaultDbName
	}
	return
}
