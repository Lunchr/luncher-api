package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"github.com/deiwin/luncher-api/lunchman/interact"
	"gopkg.in/mgo.v2/bson"
)

type restaurant struct {
	actor             interact.Actor
	collection        db.Restaurants
	usersCollection   db.Users
	regionsCollection db.Regions
}

func (r restaurant) Add() {
	checkUnique := r.getRestaurantUniquenessCheck()
	checkExists := r.getRegionExistanceCheck()

	name := getInputOrExit(r.actor, "Please enter a name for the new restaurant", checkNotEmpty, checkUnique)
	address := getInputOrExit(r.actor, "Please enter the restaurant's address", checkNotEmpty)
	region := getInputOrExit(r.actor, "Please enter the region you want to register the restaurant into", checkNotEmpty, checkExists)
	fbUserID := getInputOrExit(r.actor, "Please enter the restaurant administrator's Facebook user ID", checkNotEmpty)
	fbPageID := getInputOrExit(r.actor, "Please enter the restaurant's Facebook page ID", checkNotEmpty)

	restaurantID := r.insertRestaurantAndGetID(name, address, region)
	r.insertUser(restaurantID, fbUserID, fbPageID)

	fmt.Println("Restaurant (and user) successfully added!")
}

func (r restaurant) insertRestaurantAndGetID(name, address, region string) bson.ObjectId {
	restaurant := &model.Restaurant{
		Name:    name,
		Address: address,
		Region:  region,
	}
	confirmDBInsertion(r.actor, restaurant)
	insertedRestaurants, err := r.collection.Insert(restaurant)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	restaurantID := insertedRestaurants[0].ID
	return restaurantID
}

func (r restaurant) insertUser(restaurantID bson.ObjectId, fbUserID, fbPageID string) {
	user := &model.User{
		RestaurantID:   restaurantID,
		FacebookUserID: fbUserID,
		FacebookPageID: fbPageID,
	}
	err := r.usersCollection.Insert(user)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to enter the new user to the DB while the restaurant was already inserted. Make sure to check the DB for consistency!")
		os.Exit(1)
	}
}

func (r restaurant) getRestaurantUniquenessCheck() interact.InputCheck {
	return func(i string) error {
		if exists, err := r.collection.Exists(i); err != nil {
			return err
		} else if exists {
			return errors.New("A restaurant with the same name already exists!")
		}
		return nil
	}
}
