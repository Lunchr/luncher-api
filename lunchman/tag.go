package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/deiwin/interact"
	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"gopkg.in/mgo.v2"
)

type Tag struct {
	Actor      interact.Actor
	Collection db.Tags
}

func (t Tag) Add() {
	checkUnique := t.getTagUniquenessCheck()

	name := promptOrExit(t.Actor, "Please enter a name for the new tag", checkNotEmpty, checkSingleArg, checkUnique)
	displayName := promptOrExit(t.Actor, "Please enter a display name for the new tag", checkNotEmpty, checkSingleArg)

	t.insertTag(name, displayName)

	fmt.Println("Tag successfully added!")
}

func (t Tag) Edit(name string) {
	tag, err := t.Collection.GetName(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	checkUnique := t.getTagUniquenessCheck()

	newName := promptOptionalOrExit(t.Actor, "Please enter a name for the new tag", tag.Name, checkNotEmpty, checkSingleArg, checkUnique)
	displayName := promptOptionalOrExit(t.Actor, "Please enter a display name for the new tag", tag.DisplayName, checkNotEmpty, checkSingleArg)

	t.updateTag(name, newName, displayName)

	fmt.Println("Tag successfully updated!")
}

func (t Tag) List() {
	iter := t.Collection.GetAll()
	var tag model.Tag
	fmt.Println("Listing the tags' names and display names:")
	for iter.Next(&tag) {
		fmt.Printf("%s - %s\n", tag.Name, tag.DisplayName)
	}
	if err := iter.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (t Tag) Show(name string) {
	tag, err := t.Collection.GetName(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(pretty(tag))
}

func (t Tag) updateTag(name, newName, displayName string) {
	tag := createTag(newName, displayName)
	confirmDBInsertion(t.Actor, tag)
	err := t.Collection.UpdateName(name, tag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (t Tag) insertTag(name, displayName string) {
	tag := createTag(name, displayName)
	confirmDBInsertion(t.Actor, tag)
	if err := t.Collection.Insert(tag); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createTag(name, displayName string) *model.Tag {
	return &model.Tag{
		Name:        name,
		DisplayName: displayName,
	}
}

func (t Tag) getTagUniquenessCheck() interact.InputCheck {
	return func(i string) error {
		if _, err := t.Collection.GetName(i); err != mgo.ErrNotFound {
			if err != nil {
				return err
			}
			return errors.New("A tag with the same name already exists!")
		}
		return nil
	}
}
