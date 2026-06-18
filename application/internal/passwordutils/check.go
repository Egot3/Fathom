package passwordutils

import (
	"fmt"
	"regexp"
)

var numberRegex = regexp.MustCompile(`(\d)`)
var smallerLetterRegex = regexp.MustCompile(`([a-z])`)
var biggerLetterRegex = regexp.MustCompile(`([A-Z])`)
var forbiddenRegex = regexp.MustCompile(`([\s])`)

func CheckPasswordSafety(password string) error {
	var err error = fmt.Errorf("invalid password because it: ")
	if len(password) < 8 {
		err = fmt.Errorf("%w\n  - must contain at least 8 characters", err)
	}
	if len(numberRegex.FindStringSubmatch(password)) <= 6 { // min 6 numbers (one is always a mathced string)
		err = fmt.Errorf("%w\n  - must contain at least 6 characters", err)
	}
	if len(smallerLetterRegex.FindStringSubmatch(password)) <= 1 {
		err = fmt.Errorf("%w\n  - must contain at least 1 lowercase letter", err)
	}
	if len(biggerLetterRegex.FindStringSubmatch(password)) <= 1 {
		err = fmt.Errorf("%w\n  - must contain at least 1 uppercase letter", err)
	}
	if len(forbiddenRegex.FindString(password)) != 0 {
		err = fmt.Errorf("%w\n  - mustn't contain any type of whitespace", err)
	}

	return err
}
