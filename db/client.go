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

var ()

type Client struct {
	Database *mgo.Database
	session  *mgo.Session
}

func NewClient() *Client {
	return new(Client)
}
func (client *Client) Connect() (err error) {
	var dbURL = getEnvOrDefaultDbURL()
	var dbName = getEnvOrDefaultDbName()
	session, err := mgo.Dial(dbURL)
	if err != nil {
		return err
	}
	client.session = session
	client.Database = session.DB(dbName)
	return
}

func (client *Client) Disconnect() {
	client.session.Close()
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
