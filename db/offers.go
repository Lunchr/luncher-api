package db

import (
	"time"

	"github.com/deiwin/praad-api/db/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Offers struct {
	c *mgo.Collection
}

func NewOffers(client *Client) *Offers {
	collection := client.database.C(model.OfferCollectionName)
	return &Offers{collection}
}

func (collection Offers) Insert(offersToInsert ...*model.Offer) (err error) {
	docs := make([]interface{}, len(offersToInsert))
	for i, offer := range offersToInsert {
		docs[i] = offer
	}
	return collection.c.Insert(docs...)
}

func (collection Offers) GetForDate(date time.Time) (offers []*model.Offer, err error) {
	err = collection.c.Find(bson.M{
		"fromTime": bson.M{
			"$gte": date,
			"$lt":  date.AddDate(0, 0, 1),
		},
	}).All(&offers)
	return
}
