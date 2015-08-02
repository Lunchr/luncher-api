package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/geo"
	"github.com/deiwin/interact"
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

	name := promptOrExit(r.Actor, "Please enter a name for the new restaurant", checkNotEmpty, checkUnique)
	address := promptOrExit(r.Actor, "Please enter the restaurant's address", checkNotEmpty)
	regionName := promptOrExit(r.Actor, "Please enter the region you want to register the restaurant into", checkNotEmpty, checkExists)
	phone := promptOrExit(r.Actor, "Please enter the phone number for the restaurant", checkNotEmpty)
	fbPageID := promptOrExit(r.Actor, "Please enter the restaurant's Facebook page ID", checkNotEmpty)
	location := r.findLocationOrExit(address, regionName)

	restaurantID := r.insertRestaurantAndGetID(name, address, regionName, location, phone, fbPageID)

	fmt.Printf("Restaurant (%v) successfully added!\n", restaurantID)
}

func (r Restaurant) Edit(idString string) {
	id := bson.ObjectIdHex(idString)
	restaurant, err := r.Collection.GetID(id)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	checkUnique := r.getRestaurantUniquenessCheck()
	checkExists := r.getRegionExistanceCheck()

	name := promptOptionalOrExit(r.Actor, "Please enter a name for restaurant", restaurant.Name, checkNotEmpty, checkUnique)
	address := promptOptionalOrExit(r.Actor, "Please enter the restaurant's address", restaurant.Address, checkNotEmpty)
	regionName := promptOptionalOrExit(r.Actor, "Please enter the region you want to register the restaurant into", restaurant.Region, checkNotEmpty, checkExists)
	phone := promptOptionalOrExit(r.Actor, "Please enter the phone number for the restaurant", restaurant.Phone, checkNotEmpty)
	fbPageID := promptOptionalOrExit(r.Actor, "Please enter the restaurant's Facebook page ID", restaurant.FacebookPageID, checkNotEmpty)
	location := r.findLocationOrExit(address, regionName)

	r.updateRestaurant(id, name, address, regionName, location, phone, fbPageID)

	fmt.Println("Restaurant successfully updated!")
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
	restaurant, err := r.Collection.GetID(bson.ObjectIdHex(id))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(pretty(restaurant))
}

func (r Restaurant) updateRestaurant(id bson.ObjectId, name, address, region string, location geo.Location, phone, fbPageID string) {
	restaurant := createRestaurant(name, address, region, location, phone, fbPageID)
	confirmDBInsertion(r.Actor, restaurant)
	err := r.Collection.UpdateID(id, restaurant)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (r Restaurant) insertRestaurantAndGetID(name, address, region string, location geo.Location, phone, fbPageID string) bson.ObjectId {
	restaurant := createRestaurant(name, address, region, location, phone, fbPageID)
	confirmDBInsertion(r.Actor, restaurant)
	insertedRestaurants, err := r.Collection.Insert(restaurant)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	restaurantID := insertedRestaurants[0].ID
	return restaurantID
}

func createRestaurant(name, address, region string, location geo.Location, phone, fbPageID string) *model.Restaurant {
	return &model.Restaurant{
		Name:           name,
		Address:        address,
		Region:         region,
		Location:       model.NewPoint(location),
		Phone:          phone,
		FacebookPageID: fbPageID,
	}
}

func (r Restaurant) findLocationOrExit(address, regionName string) geo.Location {
	region, err := r.RegionsCollection.GetName(regionName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	location, err := r.Geocoder.CodeForRegion(address, region.CCTLD)
	if err == nil {
		return location
	}
	if err == geo.ErrorPartialMatch {
		message := fmt.Sprintf("Geocoder returned a partial match of (%.6f, %.6f). Do you want to continue using the partial match?", location.Lat, location.Lng)
		if confirmed, err := r.Actor.Confirm(message, interact.ConfirmDefaultToNo); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else if confirmed {
			return location
		}
	}
	fmt.Println(err)
	shouldPromptForCoords, err := r.Actor.Confirm("Do you want to enter the coordinates manually?", interact.ConfirmDefaultToYes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else if !shouldPromptForCoords {
		fmt.Println("Alright. Canceling.")
		os.Exit(1)
	}
	latitudeString := promptOrExit(r.Actor, "Please enter the latitude", checkNotEmpty, checkCanBeLatitude)
	latitude, _ := strconv.ParseFloat(latitudeString, 64)
	longitudeString := promptOrExit(r.Actor, "Please enter the longitude", checkNotEmpty, checkCanBeLongitude)
	longitude, _ := strconv.ParseFloat(longitudeString, 64)

	return geo.Location{
		Lat: latitude,
		Lng: longitude,
	}
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
		if _, err := u.RestaurantsCollection.GetID(id); err != nil {
			return err
		}
		return nil
	}
}

var (
	checkCanBeLatitude = func(i string) error {
		f, err := strconv.ParseFloat(i, 64)
		if err != nil {
			return err
		}
		if f < -90 || f > 90 {
			return errors.New("Latitude must be between -90 and 90 degrees!")
		}
		return nil
	}
	checkCanBeLongitude = func(i string) error {
		f, err := strconv.ParseFloat(i, 64)
		if err != nil {
			return err
		}
		if f < -180 || f > 180 {
			return errors.New("Longitude must be between -180 and 180 degrees!")
		}
		return nil
	}
)
