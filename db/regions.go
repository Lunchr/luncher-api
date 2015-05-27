package db

import (
	"github.com/Lunchr/luncher-api/db/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Regions interface {
	Insert(...*model.Region) error
	GetName(string) (*model.Region, error)
	GetAll() RegionIter
	UpdateName(string, *model.Region) error
}

// RegionIter is a wrapper around *mgo.Iter that allows type safe iteration
type RegionIter interface {
	Close() error
	Next(*model.Region) bool
}

type regionsCollection struct {
	*mgo.Collection
}

func NewRegions(client *Client) Regions {
	collection := client.database.C(model.RegionCollectionName)
	return &regionsCollection{collection}
}

func (c regionsCollection) Insert(regionsToInsert ...*model.Region) error {
	docs := make([]interface{}, len(regionsToInsert))
	for i, region := range regionsToInsert {
		docs[i] = region
	}
	return c.Collection.Insert(docs...)
}

func (c regionsCollection) GetName(name string) (*model.Region, error) {
	var region model.Region
	err := c.Find(bson.M{
		"name": name,
	}).One(&region)
	return &region, err
}

func (c regionsCollection) GetAll() RegionIter {
	i := c.Find(nil).Iter()
	return &regionIter{i}
}

func (c regionsCollection) UpdateName(name string, region *model.Region) error {
	return c.Update(bson.M{"name": name}, bson.M{"$set": region})
}

type regionIter struct {
	*mgo.Iter
}

func (u *regionIter) Next(region *model.Region) bool {
	return u.Iter.Next(region)
}
