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
		Ingredients []string        `json:"ingredients"          bson:"ingredients"`
		Price       float64         `json:"price"                bson:"price"`
		Tags        []string        `json:"tags"                 bson:"tags"`
		Image       string          `json:"image,omitempty"      bson:"image,omitempty"`
		FBPostID    string          `json:"fb_post_id,omitempty" bson:"fb_post_id,omitempty"`
	}

	// OfferRestaurant holds the information about the restaurant that gets included
	// in every offer
	OfferRestaurant struct {
		Name     string   `json:"name"     bson:"name"`
		Region   string   `json:"region"   bson:"region"`
		Address  string   `json:"address"  bson:"address"`
		Location Location `json:"location" bson:"location"`
	}

	// OfferWithDistance wraps an offer and adds a distance field. This struct can
	// be used to respond to queries about nearby offers.
	OfferWithDistance struct {
		Offer
		Distance float64 `json:"distance" bson:"distance"`
	}
)
