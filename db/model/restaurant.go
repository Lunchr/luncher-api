package model

import (
	"gopkg.in/mgo.v2/bson"
)

const RestaurantCollectionName = "restaurants"

type (
	Restaurant struct {
		ID      bson.ObjectId `json:"_id"           bson:"_id"`
		Name    string        `json:"name"          bson:"name"`
		Address string        `json:"address"       bson:"address"`
	}
)
