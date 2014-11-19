package db

import (
	"github.com/deiwin/praad-api/db/model"
	"gopkg.in/mgo.v2"
)

type Offers struct {
	c *mgo.Collection
}

func NewOffers(client *Client) *Offers {
	collection := client.database.C(model.OfferCollectionName)
	return &Offers{collection}
}

func (offers *Offers) Insert(offersToInsert ...*model.Offer) (err error) {
	return offers.c.Insert(offersToInsert)
}
