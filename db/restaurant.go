package db

import (
	"github.com/deiwin/luncher-api/db/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Restaurants interface {
	Insert(...*model.Restaurant) ([]*model.Restaurant, error)
	Get() ([]*model.Restaurant, error)
	GetByID(bson.ObjectId) (*model.Restaurant, error)
	Exists(name string) (bool, error)
}

type restaurantsCollection struct {
	c *mgo.Collection
}

func NewRestaurants(client *Client) Restaurants {
	collection := client.database.C(model.RestaurantCollectionName)
	return &restaurantsCollection{collection}
}

func (c restaurantsCollection) Insert(restaurantsToInsert ...*model.Restaurant) ([]*model.Restaurant, error) {
	for _, restaurant := range restaurantsToInsert {
		if restaurant.ID == "" {
			restaurant.ID = bson.NewObjectId()
		}
	}
	docs := make([]interface{}, len(restaurantsToInsert))
	for i, restaurant := range restaurantsToInsert {
		docs[i] = restaurant
	}
	return restaurantsToInsert, c.c.Insert(docs...)
}

func (c restaurantsCollection) Get() (restaurants []*model.Restaurant, err error) {
	err = c.c.Find(bson.M{}).All(&restaurants)
	return
}

func (c restaurantsCollection) GetByID(id bson.ObjectId) (*model.Restaurant, error) {
	var restaurant *model.Restaurant
	err := c.c.FindId(id).One(&restaurant)
	return restaurant, err
}

func (c restaurantsCollection) Exists(name string) (bool, error) {
	count, err := c.c.Find(bson.M{"name": name}).Count()
	if err != nil {
		return false, err
	}
	return count != 0, nil
}
