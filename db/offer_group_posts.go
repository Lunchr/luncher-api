package db

import (
	"github.com/Lunchr/luncher-api/db/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type OfferGroupPosts interface {
	Insert(...*model.OfferGroupPost) ([]*model.OfferGroupPost, error)
	UpdateByID(bson.ObjectId, *model.OfferGroupPost) error
	GetByID(bson.ObjectId) (*model.OfferGroupPost, error)
	GetByDate(model.DateWithoutTime, bson.ObjectId) (*model.OfferGroupPost, error)
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

func (c offerGroupPostCollection) UpdateByID(id bson.ObjectId, post *model.OfferGroupPost) error {
	return c.Collection.UpdateId(id, bson.M{"$set": post})
}

func (c offerGroupPostCollection) GetByID(id bson.ObjectId) (*model.OfferGroupPost, error) {
	var post model.OfferGroupPost
	err := c.FindId(id).One(&post)
	return &post, err
}

func (c offerGroupPostCollection) GetByDate(date model.DateWithoutTime, restaurantID bson.ObjectId) (*model.OfferGroupPost, error) {
	var post model.OfferGroupPost
	err := c.Find(bson.M{
		"date":          date,
		"restaurant_id": restaurantID,
	}).One(&post)
	return &post, err
}
