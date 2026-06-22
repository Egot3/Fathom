package carefulness

import (
	"errors"
	"fmt"
)

type PartialSuccess struct {
	Target      string
	WantCount   int
	ActualCount int
}

var ErrPartialSuccess = errors.New("Conflict")

func (e PartialSuccess) Error() string {
	return fmt.Sprintf("partial success on %v: %v/%v", e.Target, e.ActualCount, e.WantCount)
}

func (e PartialSuccess) JSONError() JSONError {
	return JSONError{Error: e.Error()}
}

func (e PartialSuccess) Is(target error) bool {
	return ErrConflict == target
}
