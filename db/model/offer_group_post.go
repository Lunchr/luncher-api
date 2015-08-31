package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const OfferGroupPostCollectionName = "offer_group_post"

type (
	OfferGroupPost struct {
		ID           bson.ObjectId   `json:"_id,omitempty"        bson:"_id,omitempty"`
		RestaurantID bson.ObjectId   `json:"restaurant_id"        bson:"restaurant_id"`
		Date         DateWithoutTime `json:"date"                 bson:"date"`

		MessageTemplate string `json:"message_template"     bson:"message_template"`
		FBPostID        string `json:"fb_post_id,omitempty" bson:"fb_post_id"`
	}

	DateWithoutTime string
)

const dateWithoutTimeLayout = "2006-01-02"

func DateFromTime(t time.Time) DateWithoutTime {
	dateString := t.Format(dateWithoutTimeLayout)
	return DateWithoutTime(dateString)
}

func (d DateWithoutTime) IsValid() bool {
	_, err := time.Parse(dateWithoutTimeLayout, string(d))
	return err == nil
}

func (d DateWithoutTime) TimeBounds(location *time.Location) (time.Time, time.Time, error) {
	startTime, err := time.ParseInLocation(dateWithoutTimeLayout, string(d), location)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	endTime := startTime.AddDate(0, 0, 1)
	return startTime, endTime, nil
}
