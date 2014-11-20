package db

import (
	"github.com/deiwin/praad-api/db/model"
	"gopkg.in/mgo.v2"
)

type Tags struct {
	c *mgo.Collection
}

func NewTags(client *Client) *Tags {
	collection := client.database.C(model.TagCollectionName)
	return &Tags{collection}
}

func (tags *Tags) Insert(tagsToInsert ...*model.Tag) (err error) {
	return tags.c.Insert(tagsToInsert)
}
