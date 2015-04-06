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
		Key: []string{"$geoHaystack:restaurant.location", "title"},
		// For reference, Tartu is about 0.04 lat and 0.08 lng
		// XXX BucketSize is currently missing from mgo
		// this needs to be merged: https://github.com/go-mgo/mgo/pull/90
		// BucketSize: 0.1,
	})
}
