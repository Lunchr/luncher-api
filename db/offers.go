package db

import (
	"time"

	"github.com/deiwin/luncher-api/db/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Offers interface {
	Insert(...*model.Offer) ([]*model.Offer, error)
	GetForRegion(region string, startTime, endTime time.Time) ([]*model.Offer, error)
	GetForRestaurant(restaurantName string, startTime time.Time) ([]*model.Offer, error)
	UpdateID(bson.ObjectId, *model.Offer) error
	GetID(bson.ObjectId) (*model.Offer, error)
}

type offersCollection struct {
	*mgo.Collection
}

func NewOffers(client *Client) (Offers, error) {
	collection := client.database.C(model.OfferCollectionName)
	offers := &offersCollection{collection}
	if err := offers.ensureOffersTTLIndex(); err != nil {
		return nil, err
	}
	if err := offers.ensureOffersHaystackIndex(); err != nil {
		return nil, err
	}
	return offers, nil
}

func (c offersCollection) Insert(offersToInsert ...*model.Offer) ([]*model.Offer, error) {
	for _, offer := range offersToInsert {
		if offer.ID == "" {
			offer.ID = bson.NewObjectId()
		}
	}
	docs := make([]interface{}, len(offersToInsert))
	for i, offer := range offersToInsert {
		docs[i] = offer
	}
	return offersToInsert, c.Collection.Insert(docs...)
}

func (c offersCollection) UpdateID(id bson.ObjectId, offer *model.Offer) error {
	return c.Collection.UpdateId(id, bson.M{"$set": offer})
}

func (c offersCollection) GetForRegion(region string, startTime, endTime time.Time) ([]*model.Offer, error) {
	var offers []*model.Offer
	err := c.Find(bson.M{
		"from_time": bson.M{
			"$lte": endTime,
		},
		"to_time": bson.M{
			"$gte": startTime,
		},
		"restaurant.region": region,
	}).All(&offers)
	return offers, err
}

func (c offersCollection) GetForRestaurant(restaurantName string, startTime time.Time) ([]*model.Offer, error) {
	var offers []*model.Offer
	err := c.Find(bson.M{
		"to_time": bson.M{
			"$gte": startTime,
		},
		"restaurant.name": restaurantName,
	}).All(&offers)
	return offers, err
}

func (c offersCollection) GetID(id bson.ObjectId) (*model.Offer, error) {
	var offer *model.Offer
	err := c.FindId(id).One(&offer)
	return offer, err
}
