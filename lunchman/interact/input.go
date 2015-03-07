package interact

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errCanceled = errors.New("Command aborted")
)

// InputCheck specifies the function signature for an input check
type InputCheck func(string) error

// GetInputAndRetry asks the user for input and performs the list of added checks
// on the provided input. If any of the checks fail to pass the error will be
// displayed to the user and they will then be asked if they want to try again.
// If the user does not want to retry the program will return an error.
func (a Actor) GetInputAndRetry(message string, checks ...InputCheck) (string, error) {
	for {
		input, err := a.GetInput(message, checks...)
		if err != nil {
			retryMessage := fmt.Sprintf("%v\nDo you want to try again?", err)
			confirmed, err := a.Confirm(retryMessage, ConfirmDefaultToNo)
			if err != nil {
				return "", err
			} else if !confirmed {
				return "", errCanceled
			}
			continue
		}
		return input, nil
	}
}

// GetInput asks the user for input and performs the list of added checks on the
// provided input. If any of the checks fail, the error will be returned.
func (a Actor) GetInput(message string, checks ...InputCheck) (string, error) {
	fmt.Fprint(a.w, message+": ")
	line, err := a.rd.ReadString('\n')
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
