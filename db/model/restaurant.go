package model

import (
	"labix.org/v2/mgo/bson"
)

type (
	Restaurant struct {
		ID      bson.ObjectId `json:"_id"           bson:"_id"`
		Name    string        `json:"name"          bson:"name"`
		Address string        `json:"address"       bson:"address"`
	}
)
