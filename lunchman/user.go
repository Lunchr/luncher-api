package main

import (
	"fmt"
	"os"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"github.com/deiwin/luncher-api/lunchman/interact"
	"gopkg.in/mgo.v2/bson"
)

type user struct {
	actor                 interact.Actor
	collection            db.Users
	restaurantsCollection db.Restaurants
}

func (u user) Add() {
	checkExists := u.getRestaurantExistanceCheck()

	restaurantIDString := getInputOrExit(u.actor, "Please enter the restaurant's ID this user will administrate", checkNotEmpty, checkIsObjectID, checkExists)
	restaurantID := bson.ObjectIdHex(restaurantIDString)
	fbUserID := getInputOrExit(u.actor, "Please enter the restaurant administrator's Facebook user ID", checkNotEmpty)
	fbPageID := getInputOrExit(u.actor, "Please enter the restaurant's Facebook page ID", checkNotEmpty)

	u.insertUser(restaurantID, fbUserID, fbPageID)

	fmt.Println("User successfully added!")
}

func (u user) insertUser(restaurantID bson.ObjectId, fbUserID, fbPageID string) {
	user := &model.User{
		RestaurantID:   restaurantID,
		FacebookUserID: fbUserID,
		FacebookPageID: fbPageID,
	}
	confirmDBInsertion(u.actor, user)
	err := u.collection.Insert(user)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
