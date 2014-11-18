package model

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type (
	// Offer provides the mapping to the offers as represented in the DB and also
	// to json
	Offer struct {
		ID         bson.ObjectId `json:"_id"           bson:"_id"`
		Restaurant struct {
			Name string `json:"name" bson:"name"`
		} `json:"restaurant"    bson:"restaurant"`
		Title       string    `json:"title"         bson:"title"`
		FromTime    time.Time `json:"fromTime"      bson:"fromTime"`
		ToTime      time.Time `json:"toTime"        bson:"toTime"`
		Description string    `json:"description"   bson:"description"`
		Price       int       `json:"price"         bson:"price"`
		Tags        []string  `json:"tags"          bson:"tags"`
	}
)
