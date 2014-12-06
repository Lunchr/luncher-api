package db

import (
	"github.com/deiwin/luncher-api/db/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Tags interface {
	Insert(...*model.Tag) error
	Get() ([]*model.Tag, error)
}

type tagsCollection struct {
	c *mgo.Collection
}

func NewTags(client *Client) Tags {
	collection := client.database.C(model.TagCollectionName)
	return &tagsCollection{collection}
}

func (collection tagsCollection) Insert(tagsToInsert ...*model.Tag) (err error) {
	docs := make([]interface{}, len(tagsToInsert))
	for i, tag := range tagsToInsert {
		docs[i] = tag
	}
	return collection.c.Insert(docs...)
}

func (collection tagsCollection) Get() (tags []*model.Tag, err error) {
	err = collection.c.Find(bson.M{}).All(&tags)
	return
}
