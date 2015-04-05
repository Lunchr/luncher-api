package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/deiwin/interact"
	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"github.com/deiwin/luncher-api/geo"
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

func (r Restaurant) List() {
	iter := r.Collection.GetAll()
	var restaurant model.Restaurant
	fmt.Println("Listing the restaurants' IDs and names:")
	for iter.Next(&restaurant) {
		fmt.Printf("%s - %s\n", restaurant.ID.Hex(), restaurant.Name)
	}
	if err := iter.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (r Restaurant) Show(id string) {
	restaurant, err := r.Collection.GetByID(bson.ObjectIdHex(id))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(pretty(restaurant))
}

func (r Restaurant) Edit(idString string) {
	id := bson.ObjectIdHex(idString)
	restaurant, err := r.Collection.GetByID(id)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	checkUnique := r.getRestaurantUniquenessCheck()
	checkExists := r.getRegionExistanceCheck()

	name := getInputWithDefaultOrExit(r.Actor, "Please enter a name for restaurant", restaurant.Name, checkNotEmpty, checkUnique)
	address := getInputWithDefaultOrExit(r.Actor, "Please enter the restaurant's address", restaurant.Address, checkNotEmpty)
	regionName := getInputWithDefaultOrExit(r.Actor, "Please enter the region you want to register the restaurant into", restaurant.Region, checkNotEmpty, checkExists)
	location := r.findLocationOrExit(address, regionName)

	r.updateRestaurant(id, name, address, regionName, location)

	fmt.Println("Restaurant successfully updated!")
}

func (r Restaurant) updateRestaurant(id bson.ObjectId, name, address, region string, location geo.Location) {
	restaurant := createRestaurant(name, address, region, location)
	confirmDBInsertion(r.Actor, restaurant)
	err := r.Collection.UpdateID(id, restaurant)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (r Restaurant) insertRestaurantAndGetID(name, address, region string, location geo.Location) bson.ObjectId {
	restaurant := createRestaurant(name, address, region, location)
	confirmDBInsertion(r.Actor, restaurant)
	insertedRestaurants, err := r.Collection.Insert(restaurant)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	restaurantID := insertedRestaurants[0].ID
	return restaurantID
}

func createRestaurant(name, address, region string, location geo.Location) *model.Restaurant {
	return &model.Restaurant{
		Name:     name,
		Address:  address,
		Region:   region,
		Location: location,
	}
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
