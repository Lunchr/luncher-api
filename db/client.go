package db

import "gopkg.in/mgo.v2"

type Client struct {
	config   *Config
	database *mgo.Database
	session  *mgo.Session
}

func NewClient(config *Config) *Client {
	return &Client{config: config}
}
func (client *Client) Connect() (err error) {
	session, err := mgo.Dial(client.config.DbURL)
	if err != nil {
		return err
	}
	client.session = session
	client.database = session.DB(client.config.DbName)
	return
}

func (client *Client) Disconnect() {
	client.session.Close()
}
