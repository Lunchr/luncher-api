package model

import (
	"gopkg.in/mgo.v2/bson"
)

const TagCollectionName = "tags"

type (
	Tag struct {
		ID          bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
		Name        string        `json:"name"          bson:"name"`
		DisplayName string        `json:"display_name"  bson:"display_name"`
	}
)
