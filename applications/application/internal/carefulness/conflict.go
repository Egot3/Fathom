package carefulness

import (
	"errors"
	"fmt"
	"net/http"
)

type Conflict struct {
	Conflictor string
}

var ErrConflict = errors.New("Conflict")

func (e Conflict) Error() string {
	return fmt.Sprintf("conflict on: %v", e.Conflictor)
}

func (e Conflict) JSONError() JSONError {
	return JSONError{Err: e.Error(), Status: http.StatusConflict}
}

func (e Conflict) Is(target error) bool {
	return ErrConflict == target
}

func (e Conflict) Encode(w http.ResponseWriter) {
	e.JSONError().Encode(w)
}
