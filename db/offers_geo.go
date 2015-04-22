package db

import (
	"errors"
	"time"

	"github.com/deiwin/luncher-api/db/model"
	"github.com/deiwin/luncher-api/geo"
	"gopkg.in/mgo.v2/bson"
)

func (c offersCollection) GetNear(loc geo.Location, startTime, endTime time.Time) ([]*model.OfferWithDistance, error) {
	return c.geoNear(loc, bson.M{
		"query": bson.M{
			"from_time": bson.M{
				"$lte": endTime,
			},
			"to_time": bson.M{
				"$gte": startTime,
			},
		},
		"maxDistance": 5000,
		"spherical":   true,
	})
}

func (c offersCollection) geoNear(loc geo.Location, additionalOptions bson.M) ([]*model.OfferWithDistance, error) {
	cmd := bson.D{
		{"geoNear", model.OfferCollectionName},
		{"near", model.NewPoint(loc)},
	}
	for k, v := range additionalOptions {
		cmd = append(cmd, bson.DocElem{
			Name:  k,
			Value: v,
		})
	}
	var response geoNearResponse
	if err := c.database.Run(cmd, &response); err != nil {
		return nil, err
	} else if !response.Ok {
		return nil, errors.New("The mongo geoNear command returned an error")
	}
	var offers = make([]*model.OfferWithDistance, len(response.Results))
	for i, result := range response.Results {
		offers[i] = &model.OfferWithDistance{
			Offer:    result.Obj,
			Distance: result.Dis,
		}
	}
	return offers, nil
}

type (
	geoNearResponse struct {
		Results []geoNearResult
		Ok      bool
	}

	geoNearResult struct {
		Dis float64
		Obj model.Offer
	}
)
