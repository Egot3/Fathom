package passwordutils

import (
	"errors"
	"fmt"
	"regexp"
)

var numberRegex = regexp.MustCompile(`(\d)`)
var smallerLetterRegex = regexp.MustCompile(`([a-z])`)
var biggerLetterRegex = regexp.MustCompile(`([A-Z])`)
var forbiddenRegex = regexp.MustCompile(`([\s])`)

var startingError = errors.New("invalid password because it: ")

func CheckPasswordSafety(password string) error {
	var err error = startingError
	if len(password) < 8 {
		err = fmt.Errorf("%w\n  - must contain at least 8 characters", err)
	}
	if len(password) > 60 {
		err = fmt.Errorf("%w\n  - mustn't be longer than 60 characters", err)
	}
	if len(numberRegex.FindAllString(password, -1)) < 6 { // min 6 numbers (one is always a mathced string)
		err = fmt.Errorf("%w\n  - must contain at least 6 numbers", err)
	}
	if len(smallerLetterRegex.FindAllString(password, -1)) < 1 {
		err = fmt.Errorf("%w\n  - must contain at least 1 lowercase letter", err)
	}
	if len(biggerLetterRegex.FindAllString(password, -1)) < 1 {
		err = fmt.Errorf("%w\n  - must contain at least 1 uppercase letter", err)
	}
	if len(forbiddenRegex.FindString(password)) != 0 {
		err = fmt.Errorf("%w\n  - mustn't contain any type of whitespace", err)
	}

	if err == startingError {
		return nil
	}

	return err
}
