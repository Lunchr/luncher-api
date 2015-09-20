package db

import (
	"github.com/Lunchr/luncher-api/db/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Restaurants interface {
	Insert(...*model.Restaurant) ([]*model.Restaurant, error)
	GetAll() RestaurantIter
	GetByIDs([]bson.ObjectId) ([]*model.Restaurant, error)
	GetByFacebookPageIDs([]string) ([]*model.Restaurant, error)
	GetID(bson.ObjectId) (*model.Restaurant, error)
	Exists(name string) (bool, error)
	UpdateID(bson.ObjectId, *model.Restaurant) error
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

func (c restaurantsCollection) GetAll() RestaurantIter {
	i := c.Find(nil).Iter()
	return &restaurantIter{i}
}

func (c restaurantsCollection) GetByIDs(ids []bson.ObjectId) ([]*model.Restaurant, error) {
	var restaurants []*model.Restaurant
	err := c.FindId(bson.M{
		"$in": ids,
	}).All(&restaurants)
	return restaurants, err
}

func (c restaurantsCollection) GetByFacebookPageIDs(ids []string) ([]*model.Restaurant, error) {
	var restaurants []*model.Restaurant
	err := c.Find(bson.M{
		"facebook_page_id": bson.M{
			"$in": ids,
		},
	}).All(&restaurants)
	return restaurants, err
}

func (c restaurantsCollection) GetID(id bson.ObjectId) (*model.Restaurant, error) {
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

func (c restaurantsCollection) UpdateID(id bson.ObjectId, restaurant *model.Restaurant) error {
	return c.UpdateId(id, bson.M{"$set": restaurant})
}

type restaurantIter struct {
	*mgo.Iter
}

func (u *restaurantIter) Next(restaurant *model.Restaurant) bool {
	return u.Iter.Next(restaurant)
}
