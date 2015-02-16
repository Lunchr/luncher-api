package db

import (
	"time"

	"github.com/deiwin/luncher-api/db/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Offers interface {
	Insert(...*model.Offer) error
	GetForTimeRange(time.Time, time.Time) ([]*model.Offer, error)
}

type offersCollection struct {
	c *mgo.Collection
}

func NewOffers(client *Client) Offers {
	collection := client.database.C(model.OfferCollectionName)
	return &offersCollection{collection}
}

func (collection offersCollection) Insert(offersToInsert ...*model.Offer) (err error) {
	docs := make([]interface{}, len(offersToInsert))
	for i, offer := range offersToInsert {
		docs[i] = offer
	}
	return collection.c.Insert(docs...)
}

func (collection offersCollection) GetForTimeRange(startTime time.Time, endTime time.Time) (offers []*model.Offer, err error) {
	err = collection.c.Find(bson.M{
		"from_time": bson.M{
			"$lte": endTime,
		},
		"to_time": bson.M{
			"$gte": startTime,
		},
	}).All(&offers)
	return
}
