package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"github.com/deiwin/luncher-api/lunchman/interact"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Actor                 interact.Actor
	Collection            db.Users
	RestaurantsCollection db.Restaurants
}

func (u User) Add() {
	checkExists := u.getRestaurantExistanceCheck()

	restaurantIDString := getInputOrExit(u.Actor, "Please enter the restaurant's ID this user will administrate", checkNotEmpty, checkIsObjectID, checkExists)
	restaurantID := bson.ObjectIdHex(restaurantIDString)
	fbUserID := getInputOrExit(u.Actor, "Please enter the restaurant administrator's Facebook user ID", checkNotEmpty)
	fbPageID := getInputOrExit(u.Actor, "Please enter the restaurant's Facebook page ID", checkNotEmpty)

	u.insertUser(restaurantID, fbUserID, fbPageID)

	fmt.Println("User successfully added!")
}

func (u User) List() {
	iter := u.Collection.GetAll()
	var user model.User
	for iter.Next(&user) {
		fmt.Printf("Facebook user ID: %s\n", user.FacebookUserID)
	}
	if err := iter.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (u User) Show(fbUserID string) {
	user, err := u.Collection.Get(fbUserID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	b, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(b))
}

func (u User) insertUser(restaurantID bson.ObjectId, fbUserID, fbPageID string) {
	user := &model.User{
		RestaurantID:   restaurantID,
		FacebookUserID: fbUserID,
		FacebookPageID: fbPageID,
	}
	confirmDBInsertion(u.Actor, user)
	err := u.Collection.Insert(user)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
