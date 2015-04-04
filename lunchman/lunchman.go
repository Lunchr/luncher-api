package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/lunchman/interact"
	"gopkg.in/alecthomas/kingpin.v1"
	"gopkg.in/mgo.v2/bson"
)

var (
	lunchman      = kingpin.New("lunchman", "An administrative tool to manage your luncher instance")
	add           = lunchman.Command("add", "Add a new value to the DB")
	addRegion     = add.Command("region", "Add a region")
	addRestaurant = add.Command("restaurant", "Add a restarant")
	addUser       = add.Command("user", "Add a user")

	checkNotEmpty = func(i string) error {
		if i == "" {
			return errors.New("Can't be empty!")
		}
		return nil
	}
	checkSingleArg = func(i string) error {
		if strings.Contains(i, " ") {
			return errors.New("Expecting a single argument")
		}
		return nil
	}
	checkIsObjectID = func(i string) error {
		if !bson.IsObjectIdHex(i) {
			return fmt.Errorf("%s should be, but is not an bson.ObjectId", i)
		}
		return nil
	}
)

func main() {
	dbConfig := db.NewConfig()
	dbClient := db.NewClient(dbConfig)
	err := dbClient.Connect()
	if err != nil {
		panic(err)
	}
	defer dbClient.Disconnect()

	actor := interact.NewActor(os.Stdin, os.Stdout)

	switch kingpin.MustParse(lunchman.Parse(os.Args[1:])) {
	case addRegion.FullCommand():
		region := initRegion(actor, dbClient)
		region.Add()
	case addRestaurant.FullCommand():
		restaurant := initRestaurant(actor, dbClient)
		restaurant.Add()
	case addUser.FullCommand():
		user := initUser(actor, dbClient)
		user.Add()
	}
}

func initRegion(actor interact.Actor, dbClient *db.Client) Region {
	regionsCollection := db.NewRegions(dbClient)
	return Region{actor, regionsCollection}
}

func initRestaurant(actor interact.Actor, dbClient *db.Client) Restaurant {
	restaurantsCollection := db.NewRestaurants(dbClient)
	regionsCollection := db.NewRegions(dbClient)
	return Restaurant{actor, restaurantsCollection, regionsCollection}
}

func initUser(actor interact.Actor, dbClient *db.Client) User {
	usersCollection := db.NewUsers(dbClient)
	restaurantsCollection := db.NewRestaurants(dbClient)
	return User{actor, usersCollection, restaurantsCollection}
}

func confirmDBInsertion(actor interact.Actor, o interface{}) {
	confirmationMessage := fmt.Sprintf("Going to enter the following into the DB:\n%+v\nAre you sure you want to continue?", o)
	confirmed, err := actor.Confirm(confirmationMessage, interact.ConfirmDefaultToYes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else if !confirmed {
		fmt.Println("Aborted")
		os.Exit(1)
	}
}

func getInputOrExit(a interact.Actor, message string, checks ...interact.InputCheck) string {
	input, err := a.GetInputAndRetry(message, checks...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return input
}
