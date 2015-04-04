package model

import (
	"github.com/deiwin/luncher-api/geo"
	"gopkg.in/mgo.v2/bson"
)

const RestaurantCollectionName = "restaurants"

type (
	Restaurant struct {
		ID       bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
		Name     string        `json:"name"          bson:"name"`
		Region   string        `json:"region"        bson:"region"`
		Address  string        `json:"address"       bson:"address"`
		Location geo.Location  `json:"location"      bson:"location"`
	}
)
