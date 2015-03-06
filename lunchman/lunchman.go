package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	lunchman      = kingpin.New("lunchman", "An administrative tool to manage your luncher instance")
	add           = lunchman.Command("add", "Add a new value to the DB")
	addRegion     = add.Command("region", "Add a region")
	addRestaurant = add.Command("restaurant", "Add a restarant")

	errCanceled = errors.New("Command aborted")

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
	switch kingpin.MustParse(lunchman.Parse(os.Args[1:])) {
	case addRegion.FullCommand():
		name := getInputAndRetry("Please enter a name for the new region", checkNotEmpty, checkSingleArg)
		fmt.Println("name: " + name)
	case addRestaurant.FullCommand():
		fmt.Println("add restaurant!")
	}
}

func getInputAndRetry(message string, checks ...func(string) error) string {
	for {
		input, err := getInput(message, checks...)
		if err != nil {
			c, err := confirm(fmt.Sprintf("%v\nDo you want to try again?", err), confirmDefaultToNo)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else if c == no {
				fmt.Println(errCanceled)
				os.Exit(1)
			}
			continue
		}
		return input
	}
}

func getInput(message string, checks ...func(string) error) (string, error) {
	fmt.Print(message + ": ")
	rd := bufio.NewReader(os.Stdin)
	line, err := rd.ReadString('\n')
	if err != nil {
		return "", err
	}
	input := strings.TrimSpace(line)
	for _, check := range checks {
		err = check(input)
		if err != nil {
			return "", err
		}
	}
	return input, nil
}
