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
		ID             bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
		RestaurantID   bson.ObjectId `json:"restaurant_id,omitempty" bson:"restaurant_id,omitempty"`
		FacebookUserID string        `json:"facebook_user_id" bson:"facebook_user_id"`
		Session        *UserSession  `json:"session,omitempty" bson:"session,omitempty"`
	}
	// UserSession holds data about the current user session. Some of this data
	// (facebook auth tokens, for example) may persist throughout multiple client
	// sessions, however.
	UserSession struct {
		ID                string       `bson:"id,omitempty" bson:"id,omitempty"`
		FacebookUserToken oauth2.Token `bson:"facebook_user_token,omitempty" bson:"facebook_user_token,omitempty"`
		FacebookPageToken string       `bson:"facebook_page_token,omitempty" bson:"facebook_page_token,omitempty"`
	}
)
