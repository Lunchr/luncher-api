package db

import (
	"time"

	"gopkg.in/mgo.v2"
)

const aWeek = time.Hour * 24 * 7

func (c offersCollection) ensureOffersTTLIndex() error {
	return c.EnsureIndex(mgo.Index{
		Key: []string{"to_time"},
		// This could probably be set to a day or something, but just to be safe
		// let's keep all the offers for a week
		ExpireAfter: aWeek,
	})
}

func (c offersCollection) ensureOffersHaystackIndex() error {
	return c.EnsureIndex(mgo.Index{
		Key: []string{"$2dsphere:restaurant.location", "title"},
	})
}
