package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"github.com/deiwin/luncher-api/geo"
	"github.com/deiwin/luncher-api/lunchman/interact"
	"gopkg.in/mgo.v2/bson"
)

type Restaurant struct {
	Actor             interact.Actor
	Collection        db.Restaurants
	RegionsCollection db.Regions
	Geocoder          geo.Coder
}

func (r Restaurant) Add() {
	checkUnique := r.getRestaurantUniquenessCheck()
	checkExists := r.getRegionExistanceCheck()

	name := getInputOrExit(r.Actor, "Please enter a name for the new restaurant", checkNotEmpty, checkUnique)
	address := getInputOrExit(r.Actor, "Please enter the restaurant's address", checkNotEmpty)
	regionName := getInputOrExit(r.Actor, "Please enter the region you want to register the restaurant into", checkNotEmpty, checkExists)
	location := r.findLocationOrExit(address, regionName)

	restaurantID := r.insertRestaurantAndGetID(name, address, regionName, location)

	fmt.Printf("Restaurant (%v) successfully added!\n", restaurantID)
}

func (r Restaurant) insertRestaurantAndGetID(name, address, region string, location geo.Location) bson.ObjectId {
	restaurant := &model.Restaurant{
		Name:     name,
		Address:  address,
		Region:   region,
		Location: location,
	}
	confirmDBInsertion(r.Actor, restaurant)
	insertedRestaurants, err := r.Collection.Insert(restaurant)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	restaurantID := insertedRestaurants[0].ID
	return restaurantID
}

func (r Restaurant) findLocationOrExit(address, regionName string) geo.Location {
	region, err := r.RegionsCollection.Get(regionName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	location, err := r.Geocoder.CodeForRegion(address, region.CCTLD)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return location
}

func (r Restaurant) getRestaurantUniquenessCheck() interact.InputCheck {
	return func(i string) error {
		if exists, err := r.Collection.Exists(i); err != nil {
			return err
		} else if exists {
			return errors.New("A restaurant with the same name already exists!")
		}
		return nil
	}
}

func (u User) getRestaurantExistanceCheck() interact.InputCheck {
	return func(i string) error {
		id := bson.ObjectIdHex(i)
		if _, err := u.RestaurantsCollection.GetByID(id); err != nil {
			return err
		}
		return nil
	}
}
