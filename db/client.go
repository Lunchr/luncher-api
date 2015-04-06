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

func (c *Client) Connect() (err error) {
	session, err := mgo.Dial(c.config.DbURL)
	if err != nil {
		return err
	}
	c.session = session
	c.database = session.DB(c.config.DbName)
	return
}

func (c *Client) DropDb() (err error) {
	return c.database.DropDatabase()
}

func (c *Client) Disconnect() {
	c.session.Close()
}
