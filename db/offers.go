package db

import (
	"regexp"
	"time"

	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/geo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Offers interface {
	Insert(...*model.Offer) ([]*model.Offer, error)
	GetForRegion(region string, startTime, endTime time.Time) ([]*model.Offer, error)
	GetNear(loc geo.Location, startTime, endTime time.Time) ([]*model.OfferWithDistance, error)
	GetForRestaurant(restaurantID bson.ObjectId, startTime time.Time) ([]*model.Offer, error)
	GetSimilarTitlesForRestaurant(restaurantID bson.ObjectId, partialTitle string) ([]string, error)
	GetForRestaurantByTitle(restaurantID bson.ObjectId, title string) (*model.Offer, error)
	GetForRestaurantWithinTimeBounds(restaurantID bson.ObjectId, startTime, endTime time.Time) ([]*model.Offer, error)
	UpdateID(bson.ObjectId, *model.Offer) error
	GetID(bson.ObjectId) (*model.Offer, error)
	RemoveID(bson.ObjectId) error
}

type offersCollection struct {
	*mgo.Collection
	database *mgo.Database
}

func NewOffers(c *Client) (Offers, error) {
	collection := c.database.C(model.OfferCollectionName)
	offers := &offersCollection{
		Collection: collection,
		database:   c.database,
	}
	if err := offers.ensureOffersTTLIndex(); err != nil {
		return nil, err
	}
	if err := offers.ensureOffersGeoIndex(); err != nil {
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

func (c offersCollection) GetForRestaurant(restaurantID bson.ObjectId, startTime time.Time) ([]*model.Offer, error) {
	var offers []*model.Offer
	err := c.Find(bson.M{
		"to_time": bson.M{
			"$gte": startTime,
		},
		"restaurant.id": restaurantID,
	}).All(&offers)
	return offers, err
}

func (c offersCollection) GetSimilarTitlesForRestaurant(restaurantID bson.ObjectId, partialTitle string) ([]string, error) {
	var matchingTitles []string
	err := c.Find(bson.M{
		"title": bson.M{
			"$regex":   regexp.QuoteMeta(partialTitle),
			"$options": "i",
		},
		"restaurant.id": restaurantID,
	}).Distinct("title", &matchingTitles)
	return matchingTitles, err
}

func (c offersCollection) GetForRestaurantByTitle(restaurantID bson.ObjectId, title string) (*model.Offer, error) {
	var offer model.Offer
	err := c.Find(bson.M{
		"title":         title,
		"restaurant.id": restaurantID,
	}).One(&offer)
	return &offer, err
}

func (c offersCollection) GetForRestaurantWithinTimeBounds(restaurantID bson.ObjectId, startTime, endTime time.Time) ([]*model.Offer, error) {
	var offers []*model.Offer
	err := c.Find(bson.M{
		"from_time": bson.M{
			"$lte": endTime,
		},
		"to_time": bson.M{
			"$gte": startTime,
		},
		"restaurant.id": restaurantID,
	}).All(&offers)
	return offers, err
}

func (c offersCollection) GetID(id bson.ObjectId) (*model.Offer, error) {
	var offer model.Offer
	err := c.FindId(id).One(&offer)
	return &offer, err
}

func (c offersCollection) RemoveID(id bson.ObjectId) error {
	return c.RemoveId(id)
}
