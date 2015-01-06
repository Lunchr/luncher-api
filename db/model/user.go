package model

import (
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2/bson"
)

// UserCollectionName is the collection name used in the DB for users
const UserCollectionName = "users"

type (
	// User provides the mapping to the users as represented in the DB
	User struct {
		ID                bson.ObjectId `bson:"_id,omitempty"`
		RestaurantID      bson.ObjectId `bson:"restaurant_id"`
		FacebookUserID    string        `bson:"facebook_user_id"`
		FacebookPageID    string        `bson:"facebook_page_id"`
		FacebookUserToken oauth2.Token  `bson:"facebook_user_token,omitempty"`
		FacebookPageToken string        `bson:"facebook_page_token,omitempty"`
	}
)
