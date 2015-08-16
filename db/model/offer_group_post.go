package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const OfferGroupPostCollectionName = "offer_group_post"

type (
	OfferGroupPost struct {
		ID              bson.ObjectId   `json:"_id,omitempty"        bson:"_id,omitempty"`
		RestaurantID    bson.ObjectId   `json:"restaurant_id"        bson:"restaurant_id"`
		MessageTemplate string          `json:"message_template"     bson:"message_template"`
		Date            dateWithoutTime `json:"date"                 bson:"date"`
		FBPostID        string          `json:"fb_post_id,omitempty" bson:"fb_post_id,omitempty"`
	}

	dateWithoutTime string
)

const dateWithoutTimeLayout = "2006-01-02"

func DateWithoutTime(t time.Time) dateWithoutTime {
	dateString := t.Format(dateWithoutTimeLayout)
	return dateWithoutTime(dateString)
}
