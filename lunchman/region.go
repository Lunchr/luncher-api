package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"github.com/deiwin/luncher-api/lunchman/interact"
	"gopkg.in/mgo.v2"
)

type region struct {
	actor      interact.Actor
	collection db.Regions
}

func (r region) Add() {
	checkUnique := r.getRegionUniquenessCheck()

	name := getInputOrExit(r.actor, "Please enter a name for the new region", checkNotEmpty, checkSingleArg, checkUnique)
	location := getInputOrExit(r.actor, "Please enter the region's location (IANA tz)", checkNotEmpty, checkSingleArg, checkValidLocation)
	cctld := getInputOrExit(r.actor, "Please enter the region's ccTLD (country code top-level domain)", checkNotEmpty, checkSingleArg, checkIs2Letters)

	r.insertRegion(name, location, cctld)

	fmt.Println("Region successfully added!")
}

func (r region) insertRegion(name, location, cctld string) {
	region := &model.Region{
		Name:     name,
		Location: location,
		CCTLD:    cctld,
	}
	confirmDBInsertion(r.actor, region)
	if err := r.collection.Insert(region); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (r region) getRegionUniquenessCheck() interact.InputCheck {
	return func(i string) error {
		if _, err := r.collection.Get(i); err != mgo.ErrNotFound {
			if err != nil {
				return err
			}
			return errors.New("A region with the same name already exists!")
		}
		return nil
	}
}

func (r restaurant) getRegionExistanceCheck() interact.InputCheck {
	return func(i string) error {
		if _, err := r.regionsCollection.Get(i); err != nil {
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
