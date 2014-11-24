package db

import (
	"github.com/deiwin/praad-api/db/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Tags struct {
	c *mgo.Collection
}

func NewTags(client *Client) *Tags {
	collection := client.database.C(model.TagCollectionName)
	return &Tags{collection}
}

func (collection *Tags) Insert(tagsToInsert ...*model.Tag) (err error) {
	docs := make([]interface{}, len(tagsToInsert))
	for i, tag := range tagsToInsert {
		docs[i] = tag
	}
	return collection.c.Insert(docs...)
}

func (collection Tags) Get() (tags []*model.Tag, err error) {
	err = collection.c.Find(bson.M{}).All(&tags)
	return
}
