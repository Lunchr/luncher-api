package db

import (
	"github.com/deiwin/luncher-api/db/model"
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Users interface {
	Insert(...*model.User) error
	GetFbID(string) (*model.User, error)
	GetSessionID(string) (*model.User, error)
	GetAll() UserIter
	Update(string, *model.User) error
	SetAccessToken(string, oauth2.Token) error
	SetPageAccessToken(string, string) error
	SetSessionID(string, string) error
}

// UserIter is a wrapper around *mgo.Iter that allows type safe iteration
type UserIter interface {
	Close() error
	Next(*model.User) bool
}

type usersCollection struct {
	*mgo.Collection
}

func NewUsers(client *Client) Users {
	collection := client.database.C(model.UserCollectionName)
	return &usersCollection{collection}
}

func (c usersCollection) Insert(usersToInsert ...*model.User) error {
	docs := make([]interface{}, len(usersToInsert))
	for i, user := range usersToInsert {
		docs[i] = user
	}
	return c.Collection.Insert(docs...)
}

func (c usersCollection) GetFbID(facebookUserID string) (*model.User, error) {
	var user model.User
	err := c.Find(bson.M{"facebook_user_id": facebookUserID}).One(&user)
	return &user, err
}

func (c usersCollection) GetSessionID(sessionID string) (*model.User, error) {
	var user model.User
	err := c.Find(bson.M{"session.id": sessionID}).One(&user)
	return &user, err
}

func (c usersCollection) GetAll() UserIter {
	i := c.Find(nil).Iter()
	return &userIter{i}
}

func (c usersCollection) Update(facebookUserID string, user *model.User) error {
	return c.Collection.Update(bson.M{"facebook_user_id": facebookUserID}, bson.M{"$set": user})
}

func (c usersCollection) SetAccessToken(facebookUserID string, tok oauth2.Token) error {
	return c.Collection.Update(bson.M{"facebook_user_id": facebookUserID}, bson.M{
		"$set": bson.M{"session.facebook_user_token": tok},
	})
}

func (c usersCollection) SetPageAccessToken(facebookUserID string, tok string) error {
	return c.Collection.Update(bson.M{"facebook_user_id": facebookUserID}, bson.M{
		"$set": bson.M{"session.facebook_page_token": tok},
	})
}

func (c usersCollection) SetSessionID(facebookUserID string, sessionID string) error {
	return c.Collection.Update(bson.M{"facebook_user_id": facebookUserID}, bson.M{
		"$set": bson.M{"session.id": sessionID},
	})
}

type userIter struct {
	*mgo.Iter
}

func (u *userIter) Next(user *model.User) bool {
	return u.Iter.Next(user)
}
