package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/deiwin/interact"
	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	"gopkg.in/mgo.v2"
)

type Region struct {
	Actor      interact.Actor
	Collection db.Regions
}

func (r Region) Add() {
	checkUnique := r.getRegionUniquenessCheck()

	name := promptOrExit(r.Actor, "Please enter a name for the new region", checkNotEmpty, checkSingleArg, checkUnique)
	location := promptOrExit(r.Actor, "Please enter the region's location (IANA tz)", checkNotEmpty, checkSingleArg, checkValidLocation)
	cctld := promptOrExit(r.Actor, "Please enter the region's ccTLD (country code top-level domain)", checkNotEmpty, checkSingleArg, checkIs2Letters)

	r.insertRegion(name, location, cctld)

	fmt.Println("Region successfully added!")
}

func (r Region) Edit(name string) {
	region, err := r.Collection.GetName(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	checkUnique := r.getRegionUniquenessCheck()

	newName := promptOptionalOrExit(r.Actor, "Please enter a name for the new region", region.Name, checkNotEmpty, checkSingleArg, checkUnique)
	location := promptOptionalOrExit(r.Actor, "Please enter the region's location (IANA tz)", region.Location, checkNotEmpty, checkSingleArg, checkValidLocation)
	cctld := promptOptionalOrExit(r.Actor, "Please enter the region's ccTLD (country code top-level domain)", region.CCTLD, checkNotEmpty, checkSingleArg, checkIs2Letters)

	r.updateRegion(name, newName, location, cctld)

	fmt.Println("Region successfully updated!")
}

func (r Region) List() {
	iter := r.Collection.GetAll()
	var region model.Region
	fmt.Println("Listing the regions' IDs and names:")
	for iter.Next(&region) {
		fmt.Printf("%s - %s\n", region.ID.Hex(), region.Name)
	}
	if err := iter.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (r Region) Show(name string) {
	region, err := r.Collection.GetName(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(pretty(region))
}

func (r Region) updateRegion(name, newName, location, cctld string) {
	region := createRegion(newName, location, cctld)
	confirmDBInsertion(r.Actor, region)
	err := r.Collection.UpdateName(name, region)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (r Region) insertRegion(name, location, cctld string) {
	region := createRegion(name, location, cctld)
	confirmDBInsertion(r.Actor, region)
	if err := r.Collection.Insert(region); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createRegion(name, location, cctld string) *model.Region {
	return &model.Region{
		Name:     name,
		Location: location,
		CCTLD:    cctld,
	}
}

func (r Region) getRegionUniquenessCheck() interact.InputCheck {
	return func(i string) error {
		if _, err := r.Collection.GetName(i); err != mgo.ErrNotFound {
			if err != nil {
				return err
			}
			return errors.New("A region with the same name already exists!")
		}
		return nil
	}
}

func (r Restaurant) getRegionExistanceCheck() interact.InputCheck {
	return func(i string) error {
		if _, err := r.RegionsCollection.GetName(i); err != nil {
			return err
		}
		return nil
	}
}

var (
	checkValidLocation = func(i string) error {
		if i == "Local" {
			return errors.New("Can't use region 'Local'!")
		} else if _, err := time.LoadLocation(i); err != nil {
			return err
		}
		return nil
	}
	checkIs2Letters = func(i string) error {
		if len(i) != 2 {
			return errors.New("The ccTLD should be two letters long")
		}
		return nil
	}
)
