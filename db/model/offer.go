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
		ID          bson.ObjectId   `json:"_id,omitempty"        bson:"_id,omitempty"`
		Restaurant  OfferRestaurant `json:"restaurant"           bson:"restaurant"`
		Title       string          `json:"title"                bson:"title"`
		FromTime    time.Time       `json:"from_time"            bson:"from_time"`
		ToTime      time.Time       `json:"to_time"              bson:"to_time"`
		Description string          `json:"description"          bson:"description"`
		Price       float32         `json:"price"                bson:"price"`
		Tags        []string        `json:"tags"                 bson:"tags"`
		FBPostID    string          `json:"fb_post_id,omitempty" bson:"fb_post_id,omitempty"`
	}

	OfferRestaurant struct {
		Name string `json:"name" bson:"name"`
	}
)
