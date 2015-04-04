package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/lunchman/interact"
	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	lunchman      = kingpin.New("lunchman", "An administrative tool to manage your luncher instance")
	add           = lunchman.Command("add", "Add a new value to the DB")
	addRegion     = add.Command("region", "Add a region")
	addRestaurant = add.Command("restaurant", "Add a restarant")

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
		regionsCollection := db.NewRegions(dbClient)
		region := region{actor, regionsCollection}
		region.Add()
	case addRestaurant.FullCommand():
		restaurantsCollection := db.NewRestaurants(dbClient)
		usersCollection := db.NewUsers(dbClient)
		regionsCollection := db.NewRegions(dbClient)
		restaurant := restaurant{actor, restaurantsCollection, usersCollection, regionsCollection}
		restaurant.Add()
	}
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
