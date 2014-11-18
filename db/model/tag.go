package model

import (
	"labix.org/v2/mgo/bson"
)

type (
	Tag struct {
		ID          bson.ObjectId `json:"_id"           bson:"_id"`
		Name        string        `json:"name"          bson:"name"`
		DisplayName string        `json:"displayName"   bson:"displayName"`
	}
)
