package db

import (
	"github.com/deiwin/praad-api/db/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Restaurants interface {
	Insert(...*model.Restaurant) error
	Get() ([]*model.Restaurant, error)
}

type restaurantsCollection struct {
	c *mgo.Collection
}

func NewRestaurants(client *Client) Restaurants {
	collection := client.database.C(model.RestaurantCollectionName)
	return &restaurantsCollection{collection}
}

func (collection restaurantsCollection) Insert(restaurantsToInsert ...*model.Restaurant) (err error) {
	docs := make([]interface{}, len(restaurantsToInsert))
	for i, restaurant := range restaurantsToInsert {
		docs[i] = restaurant
	}
	return collection.c.Insert(docs...)
}

func (collection restaurantsCollection) Get() (restaurants []*model.Restaurant, err error) {
	err = collection.c.Find(bson.M{}).All(&restaurants)
	return
}
