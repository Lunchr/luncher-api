package db

import (
	"github.com/deiwin/luncher-api/db/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Tags interface {
	Insert(...*model.Tag) error
	GetName(string) (*model.Tag, error)
	GetAll() TagIter
	UpdateName(string, *model.Tag) error
}

// TagIter is a wrapper around *mgo.Iter that allows type safe iteration
type TagIter interface {
	Close() error
	Next(*model.Tag) bool
}

type tagsCollection struct {
	*mgo.Collection
}

func NewTags(client *Client) Tags {
	collection := client.database.C(model.TagCollectionName)
	return &tagsCollection{collection}
}

func (c tagsCollection) Insert(tagsToInsert ...*model.Tag) (err error) {
	docs := make([]interface{}, len(tagsToInsert))
	for i, tag := range tagsToInsert {
		docs[i] = tag
	}
	return c.Collection.Insert(docs...)
}

func (c tagsCollection) GetName(name string) (*model.Tag, error) {
	var tag model.Tag
	err := c.Find(bson.M{
		"name": name,
	}).One(&tag)
	return &tag, err
}

func (c tagsCollection) GetAll() TagIter {
	i := c.Find(nil).Iter()
	return &tagIter{i}
}

func (c tagsCollection) UpdateName(name string, tag *model.Tag) error {
	return c.Update(bson.M{"name": name}, bson.M{"$set": tag})
}

type tagIter struct {
	*mgo.Iter
}

func (u *tagIter) Next(tag *model.Tag) bool {
	return u.Iter.Next(tag)
}
