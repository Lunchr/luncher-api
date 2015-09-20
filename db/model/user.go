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
		ID             bson.ObjectId   `bson:"_id,omitempty"`
		RestaurantIDs  []bson.ObjectId `bson:"restaurant_ids,omitempty"`
		FacebookUserID string          `bson:"facebook_user_id"`
		Session        UserSession     `bson:"session,omitempty"`
	}
	// UserSession holds data about the current user session. Some of this data
	// (facebook auth tokens, for example) may persist throughout multiple client
	// sessions, however.
	UserSession struct {
		ID                 string              `bson:"id,omitempty"`
		FacebookUserToken  oauth2.Token        `bson:"facebook_user_token,omitempty"`
		FacebookPageTokens []FacebookPageToken `bson:"facebook_page_tokens,omitempty"`
	}

	FacebookPageToken struct {
		PageID string `bson:"page_id"`
		Token  string `bson:"token"`
	}
)
