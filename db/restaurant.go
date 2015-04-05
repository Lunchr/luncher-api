package db

import (
	"github.com/deiwin/luncher-api/db/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Restaurants interface {
	Insert(...*model.Restaurant) ([]*model.Restaurant, error)
	Get() ([]*model.Restaurant, error)
	GetAll() RestaurantIter
	GetByID(bson.ObjectId) (*model.Restaurant, error)
	Exists(name string) (bool, error)
}

// RestaurantIter is a wrapper around *mgo.Iter that allows type safe iteration
type RestaurantIter interface {
	Close() error
	Next(*model.Restaurant) bool
}

type restaurantsCollection struct {
	*mgo.Collection
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
	return restaurantsToInsert, c.Collection.Insert(docs...)
}

func (c restaurantsCollection) Get() (restaurants []*model.Restaurant, err error) {
	err = c.Find(bson.M{}).All(&restaurants)
	return
}

func (c restaurantsCollection) GetAll() RestaurantIter {
	i := c.Find(nil).Iter()
	return &restaurantIter{i}
}

func (c restaurantsCollection) GetByID(id bson.ObjectId) (*model.Restaurant, error) {
	var restaurant model.Restaurant
	err := c.FindId(id).One(&restaurant)
	return &restaurant, err
}

func (c restaurantsCollection) Exists(name string) (bool, error) {
	count, err := c.Find(bson.M{"name": name}).Count()
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

type restaurantIter struct {
	*mgo.Iter
}

func (u restaurantIter) Next(restaurant *model.Restaurant) bool {
	return u.Iter.Next(restaurant)
}
