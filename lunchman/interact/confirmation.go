package interact

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	// ErrNoOptionSelected is returned when the user has not selected either yes or no
	ErrNoOptionSelected = errors.New("Please select y/n!")
)

// ConfirmDefault specifies what an empty user input defaults to
type ConfirmDefault int

// Possible options for what an empty user input defaults to
const (
	ConfirmDefaultToYes ConfirmDefault = iota
	ConfirmDefaultToNo
	ConfirmNoDefault
)

// Confirm provides the message to the user and asks yes or no. If the user
// doesn't select either of the possible answers ErrNoOptionSelected will be
// returned
func Confirm(message string, def ConfirmDefault) (bool, error) {
	var options string
	switch def {
	case ConfirmDefaultToYes:
		options = "[Y/n]"
	case ConfirmDefaultToNo:
		options = "[y/N]"
	case ConfirmNoDefault:
		options = "[y/n]"
	}
	fmt.Printf("%s %s: ", message, options)

	rd := bufio.NewReader(os.Stdin)
	line, err := rd.ReadString('\n')
	input := strings.TrimSpace(line)
	if err != nil {
		return false, err
	} else if input == "" {
		switch def {
		case ConfirmDefaultToYes:
			return true, nil
		case ConfirmDefaultToNo:
			return false, nil
		case ConfirmNoDefault:
			return false, ErrNoOptionSelected
		}
	}
	switch input {
	case "y":
		return true, nil
	case "n":
		return false, nil
	}
	return false, ErrNoOptionSelected
}
