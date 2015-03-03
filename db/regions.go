package db

import (
	"github.com/deiwin/luncher-api/db/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Regions interface {
	Insert(...*model.Region) error
	Get(string) (*model.Region, error)
}

type regionsCollection struct {
	c *mgo.Collection
}

func NewRegions(client *Client) Regions {
	collection := client.database.C(model.RegionCollectionName)
	return &regionsCollection{collection}
}

func (collection regionsCollection) Insert(regionsToInsert ...*model.Region) (err error) {
	docs := make([]interface{}, len(regionsToInsert))
	for i, region := range regionsToInsert {
		docs[i] = region
	}
	return collection.c.Insert(docs...)
}

func (collection regionsCollection) Get(name string) (region *model.Region, err error) {
	err = collection.c.Find(bson.M{
		"name": name,
	}).One(&region)
	return
}
