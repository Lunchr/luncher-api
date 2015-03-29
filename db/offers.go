package db

import (
	"time"

	"github.com/deiwin/luncher-api/db/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Offers interface {
	Insert(...*model.Offer) ([]*model.Offer, error)
	Get(region string, startTime, endTime time.Time) ([]*model.Offer, error)
	UpdateID(bson.ObjectId, *model.Offer) error
	GetByID(bson.ObjectId) (*model.Offer, error)
}

type offersCollection struct {
	c *mgo.Collection
}

func NewOffers(client *Client) Offers {
	collection := client.database.C(model.OfferCollectionName)
	return &offersCollection{collection}
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
	return offersToInsert, c.c.Insert(docs...)
}

func (c offersCollection) UpdateID(id bson.ObjectId, offer *model.Offer) error {
	return c.c.UpdateId(id, bson.M{"$set": offer})
}

func (c offersCollection) Get(region string, startTime, endTime time.Time) (offers []*model.Offer, err error) {
	err = c.c.Find(bson.M{
		"from_time": bson.M{
			"$lte": endTime,
		},
		"to_time": bson.M{
			"$gte": startTime,
		},
		"restaurant.region": region,
	}).All(&offers)
	return
}

func (c offersCollection) GetByID(id bson.ObjectId) (*model.Offer, error) {
	var offer *model.Offer
	err := c.c.FindId(id).One(&offer)
	return offer, err
}
