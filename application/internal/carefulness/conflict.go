package carefulness

import (
	"errors"
	"fmt"
)

type Conflict struct {
	Conflictor string
}

var ErrConflict = errors.New("Conflict")

func (e Conflict) Error() string {
	return fmt.Sprintf("conflict on: %v", e.Conflictor)
}

func (e Conflict) JSONError() JSONError {
	return JSONError{Error: e.Error()}
}

func (e Conflict) Is(target error) bool {
	return ErrConflict == target
}
