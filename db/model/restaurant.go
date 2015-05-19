package model

import (
	"github.com/deiwin/luncher-api/geo"
	"gopkg.in/mgo.v2/bson"
)

const RestaurantCollectionName = "restaurants"

type (
	Restaurant struct {
		ID       bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
		Name     string        `json:"name"          bson:"name"`
		Region   string        `json:"region"        bson:"region"`
		Address  string        `json:"address"       bson:"address"`
		Location Location      `json:"location"      bson:"location"`
		Phone    string        `json:"phone"         bson:"phone"`
	}

	// Location is a (limited) representation of a GeoJSON object
	Location struct {
		Type        string    `json:"type" bson:"type"`
		Coordinates []float64 `json:"coordinates" bson:"coordinates"`
	}
)

// NewPoint creates a Location object that get's marshalled into a GeoJSON object
// which mongo recognizes
func NewPoint(loc geo.Location) Location {
	return Location{
		Type:        "Point",
		Coordinates: []float64{loc.Lng, loc.Lat},
	}
}
