package db

import (
	"time"

	"github.com/Lunchr/luncher-api/db/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type RegistrationAccessTokens interface {
	Insert(*model.RegistrationAccessToken) (*model.RegistrationAccessToken, error)
	Exists(model.Token) (bool, error)
}

type registrationAccessTokensCollection struct {
	*mgo.Collection
}

func NewRegistrationAccessTokens(c *Client) (RegistrationAccessTokens, error) {
	collection := c.database.C(model.RegistrationAccessTokenCollectionName)
	tokens := &registrationAccessTokensCollection{collection}
	if err := tokens.ensureTTLIndex(); err != nil {
		return nil, err
	}
	return tokens, nil
}

func (c registrationAccessTokensCollection) Insert(t *model.RegistrationAccessToken) (*model.RegistrationAccessToken, error) {
	if t.ID == "" {
		// TODO copy before changing perhaps. Mutating incoming pointed objects can be bad
		t.ID = bson.NewObjectId()
	}
	return t, c.Collection.Insert(t)
}

func (c registrationAccessTokensCollection) Exists(token model.Token) (bool, error) {
	count, err := c.Find(bson.M{
		"token": token,
	}).Count()
	return count > 0, err
}

func (c registrationAccessTokensCollection) ensureTTLIndex() error {
	return c.EnsureIndex(mgo.Index{
		Key:         []string{"created_at"},
		ExpireAfter: time.Hour * 24 * 7, // 7 days
	})
}
