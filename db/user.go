package db

import (
	"github.com/deiwin/luncher-api/db/model"
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Users interface {
	Insert(...*model.User) error
	Get(string) (*model.User, error)
	SetAccessToken(string, oauth2.Token) error
	SetPageAccessToken(string, string) error
}

type usersCollection struct {
	c *mgo.Collection
}

func NewUsers(client *Client) Users {
	collection := client.database.C(model.UserCollectionName)
	return &usersCollection{collection}
}

func (collection usersCollection) Insert(usersToInsert ...*model.User) error {
	docs := make([]interface{}, len(usersToInsert))
	for i, user := range usersToInsert {
		docs[i] = user
	}
	return collection.c.Insert(docs...)
}

func (collection usersCollection) Get(facebookUserID string) (user *model.User, err error) {
	err = collection.c.Find(bson.M{"facebook_user_id": facebookUserID}).One(&user)
	return
}

func (collection usersCollection) SetAccessToken(facebookUserID string, tok oauth2.Token) error {
	return collection.c.Update(bson.M{"facebook_user_id": facebookUserID}, bson.M{
		"$set": bson.M{"facebook_user_token": tok},
	})
}

func (collection usersCollection) SetPageAccessToken(facebookUserID string, tok string) error {
	return collection.c.Update(bson.M{"facebook_user_id": facebookUserID}, bson.M{
		"$set": bson.M{"facebook_page_token": tok},
	})
}
