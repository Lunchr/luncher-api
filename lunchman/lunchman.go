package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	lunchman = kingpin.New("lunchman", "An administrative tool to manage your luncher instance")
	add      = lunchman.Command("add", "Add a new value to the DB")

	addRegion = add.Command("region", "Add a region")

	addRestaurant = add.Command("restaurant", "Add a restarant")
)

func main() {
	switch kingpin.MustParse(lunchman.Parse(os.Args[1:])) {
	case addRegion.FullCommand():
		fmt.Println("add region!")
	case addRestaurant.FullCommand():
		fmt.Println("add restaurant!")
	}
}
