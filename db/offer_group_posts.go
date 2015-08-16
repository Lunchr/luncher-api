package db

import (
	"github.com/Lunchr/luncher-api/db/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type OfferGroupPosts interface {
	Insert(...*model.OfferGroupPost) ([]*model.OfferGroupPost, error)
}

type offerGroupPostCollection struct {
	*mgo.Collection
}

func NewOfferGroupPosts(c *Client) OfferGroupPosts {
	collection := c.database.C(model.OfferGroupPostCollectionName)
	return &offerGroupPostCollection{collection}
}

func (c offerGroupPostCollection) Insert(posts ...*model.OfferGroupPost) ([]*model.OfferGroupPost, error) {
	for _, offerGroupPost := range posts {
		if offerGroupPost.ID == "" {
			offerGroupPost.ID = bson.NewObjectId()
		}
	}
	docs := make([]interface{}, len(posts))
	for i, post := range posts {
		docs[i] = post
	}
	return posts, c.Collection.Insert(docs...)
}
