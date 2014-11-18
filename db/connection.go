package db

import (
	"os"

	"gopkg.in/mgo.v2"
)

const (
	dbURLEnvVariable  = "PRAAD_DB_ADDRESS"
	dbNameEnvVariable = "PRAAD_DB_NAME"
	testDbURL         = "mongodb://localhost/test"
	testDbName        = "test"
)

var (
	Database *mgo.Database
	session  *mgo.Session
)

func Connect() {
	var err error
	var dbURL = getEnvOrDefaultDbURL()
	var dbName = getEnvOrDefaultDbName()
	session, err = mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}
	Database = session.DB(dbName)
}

func Disconnect() {
	session.Close()
}

func getEnvOrDefaultDbURL() (dbURL string) {
	dbURL = os.Getenv(dbURLEnvVariable)
	if dbURL == "" {
		dbURL = testDbURL
	}
	return
}

func getEnvOrDefaultDbName() (dbName string) {
	dbName = os.Getenv(dbNameEnvVariable)
	if dbName == "" {
		dbName = testDbName
	}
	return
}
