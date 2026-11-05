package passwordutils

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Policy struct {
	MinLength       int  `json:"min_length"`
	MaxLength       int  `json:"max_length"`
	MinDigits       int  `json:"min_digits"`
	MinLowercase    int  `json:"min_lowercase"`
	MinUppercase    int  `json:"min_uppercase"`
	AllowWhitespace bool `json:"allow_whitespace"`
	MaxBytes        int  `json:"max_bytes"`
}

var DefaultPolicy = Policy{
	MinLength:    8,
	MaxLength:    60,
	MinDigits:    6,
	MinLowercase: 1,
	MinUppercase: 1,
	MaxBytes:     72,
}

func (p Policy) Violations(password string) []string {
	var digits, lower, upper int
	var hasSpace bool
	for _, r := range password {
		switch {
		case r >= '0' && r <= '9':
			digits++
		case r >= 'a' && r <= 'z':
			lower++
		case r >= 'A' && r <= 'Z':
			upper++
		}
		if unicode.IsSpace(r) {
			hasSpace = true
		}
	}

	var v []string
	if n := utf8.RuneCountInString(password); n < p.MinLength {
		v = append(v, fmt.Sprintf("must contain at least %d characters", p.MinLength))
	} else if n > p.MaxLength {
		v = append(v, fmt.Sprintf("mustn't be longer than %d characters", p.MaxLength))
	}
	if len(password) > p.MaxBytes {
		v = append(v, fmt.Sprintf("mustn't be longer than %d bytes", p.MaxBytes))
	}
	if digits < p.MinDigits {
		v = append(v, fmt.Sprintf("must contain at least %d numbers", p.MinDigits))
	}
	if lower < p.MinLowercase {
		v = append(v, fmt.Sprintf("must contain at least %d lowercase letters", p.MinLowercase))
	}
	if upper < p.MinUppercase {
		v = append(v, fmt.Sprintf("must contain at least %d uppercase letters", p.MinUppercase))
	}
	if !p.AllowWhitespace && hasSpace {
		v = append(v, "mustn't contain any type of whitespace")
	}
	return v
}

func CheckPasswordSafety(password string) error {
	v := DefaultPolicy.Violations(password)
	if len(v) == 0 {
		return nil
	}
	return fmt.Errorf("invalid password because it:\n  - %s", strings.Join(v, "\n  - "))
}
