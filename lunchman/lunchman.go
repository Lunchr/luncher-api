package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"github.com/deiwin/luncher-api/lunchman/interact"
	"gopkg.in/alecthomas/kingpin.v1"
	"gopkg.in/mgo.v2"
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
	checkValidLocation = func(i string) error {
		if i == "Local" {
			return errors.New("Can't use region 'Local'!")
		} else if _, err := time.LoadLocation(i); err != nil {
			return err
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

	actor := interact.NewActor(os.Stdin)

	switch kingpin.MustParse(lunchman.Parse(os.Args[1:])) {
	case addRegion.FullCommand():
		regionsCollection := db.NewRegions(dbClient)
		checkUnique := getRegionUniquenessCheck(regionsCollection)
		name := getInputOrExit(actor, "Please enter a name for the new region", checkNotEmpty, checkSingleArg, checkUnique)
		locInput := getInputOrExit(actor, "Please enter the region's location (IANA tz)", checkNotEmpty, checkSingleArg, checkValidLocation)
		region := &model.Region{
			Name:     name,
			Location: locInput,
		}
		if err = regionsCollection.Insert(region); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Region successfully added!")
	case addRestaurant.FullCommand():
		restaurantsCollection := db.NewRestaurants(dbClient)
		checkUnique := getRestaurantUniquenessCheck(restaurantsCollection)
		name := getInputOrExit(actor, "Please enter a name for the new region", checkNotEmpty, checkUnique)
		address := getInputOrExit(actor, "Please enter the restaurant's address", checkNotEmpty)
		// FB user ID
		// FB page ID
		fmt.Println("add restaurant!")
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

func getRegionUniquenessCheck(c db.Regions) interact.InputCheck {
	return func(i string) error {
		if _, err := c.Get(i); err != mgo.ErrNotFound {
			if err != nil {
				return err
			}
			return errors.New("A region with the same name already exists!")
		}
		return nil
	}
}

func getRestaurantUniquenessCheck(c db.Restaurants) interact.InputCheck {
	return func(i string) error {
		if exists, err := c.Exists(i); err != nil {
			return err
		} else if exists {
			return errors.New("A restaurant with the same name already exists!")
		}
		return nil
	}
}
