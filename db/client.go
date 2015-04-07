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

func (c *Client) Connect() error {
	session, err := mgo.Dial(c.config.DbURL)
	if err != nil {
		return err
	}
	c.session = session
	c.database = session.DB(c.config.DbName)
	return err
}

func (c *Client) WipeDb() error {
	c.session.ResetIndexCache()
	return c.database.DropDatabase()
}

func (c *Client) Disconnect() {
	c.session.Close()
}
