package db

import (
	"github.com/deiwin/praad-api/db/model"
	"gopkg.in/mgo.v2"
)

type Restaurants struct {
	c *mgo.Collection
}

func NewRestaurants(client *Client) *Restaurants {
	collection := client.database.C(model.RestaurantCollectionName)
	return &Restaurants{collection}
}

func (restaurants *Restaurants) Insert(restaurantsToInsert ...*model.Restaurant) (err error) {
	return restaurants.c.Insert(restaurantsToInsert)
}
