package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/deiwin/interact"
	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"gopkg.in/mgo.v2"
)

type Region struct {
	Actor      interact.Actor
	Collection db.Regions
}

func (r Region) Add() {
	checkUnique := r.getRegionUniquenessCheck()

	name := getInputOrExit(r.Actor, "Please enter a name for the new region", checkNotEmpty, checkSingleArg, checkUnique)
	location := getInputOrExit(r.Actor, "Please enter the region's location (IANA tz)", checkNotEmpty, checkSingleArg, checkValidLocation)
	cctld := getInputOrExit(r.Actor, "Please enter the region's ccTLD (country code top-level domain)", checkNotEmpty, checkSingleArg, checkIs2Letters)

	r.insertRegion(name, location, cctld)

	fmt.Println("Region successfully added!")
}

func (r Region) insertRegion(name, location, cctld string) {
	region := &model.Region{
		Name:     name,
		Location: location,
		CCTLD:    cctld,
	}
	confirmDBInsertion(r.Actor, region)
	if err := r.Collection.Insert(region); err != nil {
		fmt.Println(err)
		os.Exit(1)
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
