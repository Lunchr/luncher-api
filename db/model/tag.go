package model

import (
	"gopkg.in/mgo.v2/bson"
)

const TagCollectionName = "tags"

type (
	Tag struct {
		ID          bson.ObjectId `json:"_id"           bson:"_id"`
		Name        string        `json:"name"          bson:"name"`
		DisplayName string        `json:"displayName"   bson:"displayName"`
	}
)
