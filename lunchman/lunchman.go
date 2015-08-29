package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/geo"
	"github.com/deiwin/interact"
	"gopkg.in/alecthomas/kingpin.v1"
	"gopkg.in/mgo.v2/bson"
)

var (
	lunchman             = kingpin.New("lunchman", "An administrative tool to manage your luncher instance")
	add                  = lunchman.Command("add", "Add a new value to the DB")
	addRegion            = add.Command("region", "Add a region")
	addRestaurant        = add.Command("restaurant", "Add a restarant")
	addUser              = add.Command("user", "Add a user")
	addTag               = add.Command("tag", "Add a tag")
	addRegistrationToken = add.Command("token", "Create and add a new registration access token")

	list            = lunchman.Command("list", "List the current values in DB")
	listRegions     = list.Command("regions", "List all regions")
	listRestaurants = list.Command("restaurants", "List all restaurants")
	listUsers       = list.Command("users", "List all users")
	listTags        = list.Command("tags", "List all tags")

	show             = lunchman.Command("show", "Show a specific DB item")
	showRegion       = show.Command("region", "Show a region")
	showRegionName   = showRegion.Arg("name", "The region's name").Required().String()
	showRestaurant   = show.Command("restaurant", "Show a restaurant")
	showRestaurantID = showRestaurant.Arg("id", "The restaurant's ID").Required().String()
	showUser         = show.Command("user", "Show a user")
	showUserID       = showUser.Arg("facebookid", "The user's Facebook ID").Required().String()
	showTag          = show.Command("tag", "Show a tag")
	showTagName      = showTag.Arg("name", "The tag's name").Required().String()

	edit             = lunchman.Command("edit", "Edit a specific DB item")
	editRegion       = edit.Command("region", "Edit a region")
	editRegionName   = editRegion.Arg("name", "The region's name").Required().String()
	editRestaurant   = edit.Command("restaurant", "Edit a restaurant")
	editRestaurantID = editRestaurant.Arg("id", "The restaurant's ID").Required().String()
	editUser         = edit.Command("user", "Edit a user")
	editUserID       = editUser.Arg("facebookid", "The user's Facebook ID").Required().String()
	editTag          = edit.Command("tag", "Edit a tag")
	editTagName      = editTag.Arg("name", "The tag's name").Required().String()

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
	case addTag.FullCommand():
		tag := initTag(actor, dbClient)
		tag.Add()
	case addRegistrationToken.FullCommand():
		token := initRegistrationToken(dbClient)
		token.CreateAndAdd()

	case listRegions.FullCommand():
		region := initRegion(actor, dbClient)
		region.List()
	case listRestaurants.FullCommand():
		restaurant := initRestaurant(actor, dbClient)
		restaurant.List()
	case listUsers.FullCommand():
		user := initUser(actor, dbClient)
		user.List()
	case listTags.FullCommand():
		tag := initTag(actor, dbClient)
		tag.List()

	case showRegion.FullCommand():
		region := initRegion(actor, dbClient)
		region.Show(*showRegionName)
	case showRestaurant.FullCommand():
		restaurant := initRestaurant(actor, dbClient)
		restaurant.Show(*showRestaurantID)
	case showUser.FullCommand():
		user := initUser(actor, dbClient)
		user.Show(*showUserID)
	case showTag.FullCommand():
		tag := initTag(actor, dbClient)
		tag.Show(*showTagName)

	case editRegion.FullCommand():
		region := initRegion(actor, dbClient)
		region.Edit(*editRegionName)
	case editRestaurant.FullCommand():
		restaurant := initRestaurant(actor, dbClient)
		restaurant.Edit(*editRestaurantID)
	case editUser.FullCommand():
		user := initUser(actor, dbClient)
		user.Edit(*editUserID)
	case editTag.FullCommand():
		tag := initTag(actor, dbClient)
		tag.Edit(*editTagName)
	}
}

func initRegion(actor interact.Actor, dbClient *db.Client) Region {
	regionsCollection := db.NewRegions(dbClient)
	return Region{actor, regionsCollection}
}

func initRestaurant(actor interact.Actor, dbClient *db.Client) Restaurant {
	restaurantsCollection := db.NewRestaurants(dbClient)
	regionsCollection := db.NewRegions(dbClient)
	geoConf := geo.NewConfig()
	geocoder := geo.NewCoder(geoConf)
	return Restaurant{actor, restaurantsCollection, regionsCollection, geocoder}
}

func initUser(actor interact.Actor, dbClient *db.Client) User {
	usersCollection := db.NewUsers(dbClient)
	restaurantsCollection := db.NewRestaurants(dbClient)
	return User{actor, usersCollection, restaurantsCollection}
}

func initTag(actor interact.Actor, dbClient *db.Client) Tag {
	tagsCollection := db.NewTags(dbClient)
	return Tag{actor, tagsCollection}
}

func initRegistrationToken(dbClient *db.Client) RegistrationToken {
	collection, err := db.NewRegistrationAccessTokens(dbClient)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return RegistrationToken{collection}
}

func confirmDBInsertion(actor interact.Actor, o interface{}) {
	confirmationMessage := fmt.Sprintf("Going to enter the following into the DB:\n%s\nAre you sure you want to continue?", pretty(o))
	confirmed, err := actor.Confirm(confirmationMessage, interact.ConfirmDefaultToYes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else if !confirmed {
		fmt.Println("Aborted")
		os.Exit(1)
	}
}

func promptOrExit(a interact.Actor, message string, checks ...interact.InputCheck) string {
	input, err := a.PromptAndRetry(message, checks...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return input
}

func promptOptionalOrExit(a interact.Actor, message, fallback string, checks ...interact.InputCheck) string {
	input, err := a.PromptOptionalAndRetry(message, fallback, checks...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return input
}

func pretty(o interface{}) string {
	b, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(b)
}
