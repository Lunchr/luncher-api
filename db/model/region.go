package model

import (
	"gopkg.in/mgo.v2/bson"
)

const RegionCollectionName = "regions"

type (
	// A Region specifies a town or a district supported by the application and
	// the data related to this region, such as the time zone.
	Region struct {
		ID       bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
		Name     string        `json:"name"          bson:"name"`
		Location string        `json:"location"      bson:"location"`
	}
)
