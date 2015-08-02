package main

import (
	"fmt"
	"os"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	"github.com/deiwin/interact"
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

	u.insertUser(restaurantID, fbUserID)

	fmt.Println("User successfully added!")
}

func (u User) Edit(fbUserID string) {
	user, err := u.Collection.GetFbID(fbUserID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	checkExists := u.getRestaurantExistanceCheck()

	restaurantIDString := promptOptionalOrExit(u.Actor, "Please enter the restaurant's ID this user will administrate", user.RestaurantIDs[0].Hex(), checkNotEmpty, checkIsObjectID, checkExists)
	restaurantID := bson.ObjectIdHex(restaurantIDString)
	newFBUserID := promptOptionalOrExit(u.Actor, "Please enter the restaurant administrator's Facebook user ID", user.FacebookUserID, checkNotEmpty)

	u.updateUser(fbUserID, restaurantID, newFBUserID)

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

func (u User) updateUser(fbUserID string, restaurantID bson.ObjectId, newFBUserID string) {
	user := createUser(restaurantID, newFBUserID)
	confirmDBInsertion(u.Actor, user)
	err := u.Collection.Update(fbUserID, user)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (u User) insertUser(restaurantID bson.ObjectId, fbUserID string) {
	user := createUser(restaurantID, fbUserID)
	confirmDBInsertion(u.Actor, user)
	err := u.Collection.Insert(user)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createUser(restaurantID bson.ObjectId, fbUserID string) *model.User {
	return &model.User{
		RestaurantIDs:  []bson.ObjectId{restaurantID},
		FacebookUserID: fbUserID,
	}
}
