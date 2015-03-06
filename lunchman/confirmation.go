package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	errNoOptionSelected = errors.New("Please select y/n!")
)

type confirmDefault int

const (
	confirmDefaultToYes confirmDefault = iota
	confirmDefaultToNo
	confirmNoDefault
)

type confirmResult int

const (
	yes confirmResult = iota
	no
)

func confirm(message string, def confirmDefault) (confirmResult, error) {
	var options string
	switch def {
	case confirmDefaultToYes:
		options = "[Y/n]"
	case confirmDefaultToNo:
		options = "[y/N]"
	case confirmNoDefault:
		options = "[y/n]"
	}
	fmt.Printf("%s %s: ", message, options)

	rd := bufio.NewReader(os.Stdin)
	line, err := rd.ReadString('\n')
	input := strings.TrimSpace(line)
	if err != nil {
		return 0, err
	} else if input == "" {
		switch def {
		case confirmDefaultToYes:
			return yes, nil
		case confirmDefaultToNo:
			return no, nil
		case confirmNoDefault:
			return 0, errNoOptionSelected
		}
	}
	switch input {
	case "y":
		return yes, nil
	case "n":
		return no, nil
	}
	return 0, errNoOptionSelected
}
