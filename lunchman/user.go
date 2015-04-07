package main

import (
	"fmt"
	"os"

	"github.com/deiwin/interact"
	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Actor                 interact.Actor
	Collection            db.Users
	RestaurantsCollection db.Restaurants
}

func (u User) Add() {
	checkExists := u.getRestaurantExistanceCheck()

	restaurantIDString := promptOrExit(u.Actor, "Please enter the restaurant's ID this user will administrate", checkNotEmpty, checkIsObjectID, checkExists)
	restaurantID := bson.ObjectIdHex(restaurantIDString)
	fbUserID := promptOrExit(u.Actor, "Please enter the restaurant administrator's Facebook user ID", checkNotEmpty)
	fbPageID := promptOrExit(u.Actor, "Please enter the restaurant's Facebook page ID", checkNotEmpty)

	u.insertUser(restaurantID, fbUserID, fbPageID)

	fmt.Println("User successfully added!")
}

func (u User) Edit(fbUserID string) {
	user, err := u.Collection.GetFbID(fbUserID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	checkExists := u.getRestaurantExistanceCheck()

	restaurantIDString := promptOptionalOrExit(u.Actor, "Please enter the restaurant's ID this user will administrate", user.RestaurantID.Hex(), checkNotEmpty, checkIsObjectID, checkExists)
	restaurantID := bson.ObjectIdHex(restaurantIDString)
	newFBUserID := promptOptionalOrExit(u.Actor, "Please enter the restaurant administrator's Facebook user ID", user.FacebookUserID, checkNotEmpty)
	fbPageID := promptOptionalOrExit(u.Actor, "Please enter the restaurant's Facebook page ID", user.FacebookPageID, checkNotEmpty)

	u.updateUser(fbUserID, restaurantID, newFBUserID, fbPageID)

	fmt.Println("User successfully updated!")
}

func (u User) List() {
	iter := u.Collection.GetAll()
	var user model.User
	fmt.Println("Listing the users' Facebook user IDs:")
	for iter.Next(&user) {
		fmt.Println(user.FacebookUserID)
	}
	if err := iter.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (u User) Show(fbUserID string) {
	user, err := u.Collection.GetFbID(fbUserID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(pretty(user))
}

func (u User) updateUser(fbUserID string, restaurantID bson.ObjectId, newFBUserID, fbPageID string) {
	user := createUser(restaurantID, newFBUserID, fbPageID)
	confirmDBInsertion(u.Actor, user)
	err := u.Collection.Update(fbUserID, user)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (u User) insertUser(restaurantID bson.ObjectId, fbUserID, fbPageID string) {
	user := createUser(restaurantID, fbUserID, fbPageID)
	confirmDBInsertion(u.Actor, user)
	err := u.Collection.Insert(user)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createUser(restaurantID bson.ObjectId, fbUserID, fbPageID string) *model.User {
	return &model.User{
		RestaurantID:   restaurantID,
		FacebookUserID: fbUserID,
		FacebookPageID: fbPageID,
	}
}
