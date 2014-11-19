package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const OfferCollectionName = "offers"

type (
	// Offer provides the mapping to the offers as represented in the DB and also
	// to json
	Offer struct {
		ID          bson.ObjectId   `json:"_id,omitempty"            bson:"_id,omitempty"`
		Restaurant  OfferRestaurant `json:"restaurant"     bson:"restaurant"`
		Title       string          `json:"title"          bson:"title"`
		FromTime    time.Time       `json:"fromTime"       bson:"fromTime"`
		ToTime      time.Time       `json:"toTime"         bson:"toTime"`
		Description string          `json:"description"    bson:"description"`
		Price       float32         `json:"price"          bson:"price"`
		Tags        []string        `json:"tags"           bson:"tags"`
	}

	OfferRestaurant struct {
		Name string `json:"name" bson:"name"`
	}
)
