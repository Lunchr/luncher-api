package model

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type (
	Offer struct {
		Id          bson.ObjectId `json:"_id"           bson:"_id"`
		Restaurant  Restaurant    `json:"restaurant"    bson:"restaurant"`
		Title       string        `json:"title"         bson:"title"`
		FromTime    time.Time     `json:"fromTime"      bson:"fromTime"`
		ToTime      time.Time     `json:"toTime"        bson:"toTime"`
		Description string        `json:"description"   bson:"description"`
		Price       int           `json:"price"         bson:"price"`
		Tags        []Tag         `json:"tags"          bson:"tags"`
	}
)
