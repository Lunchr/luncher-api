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
	regionsCollection db.Regions
}

func (r restaurant) Add() {
	checkUnique := r.getRestaurantUniquenessCheck()
	checkExists := r.getRegionExistanceCheck()

	name := getInputOrExit(r.actor, "Please enter a name for the new restaurant", checkNotEmpty, checkUnique)
	address := getInputOrExit(r.actor, "Please enter the restaurant's address", checkNotEmpty)
	region := getInputOrExit(r.actor, "Please enter the region you want to register the restaurant into", checkNotEmpty, checkExists)

	restaurantID := r.insertRestaurantAndGetID(name, address, region)

	fmt.Printf("Restaurant (%v) successfully added!\n", restaurantID)
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

func (u user) getRestaurantExistanceCheck() interact.InputCheck {
	return func(i string) error {
		id := bson.ObjectIdHex(i)
		if _, err := u.restaurantsCollection.GetByID(id); err != nil {
			return err
		}
		return nil
	}
}
