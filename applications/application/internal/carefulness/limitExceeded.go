package carefulness

import (
	"errors"
	"net/http"
)

type limitExceeded struct{}

var ErrLimitExceeded *limitExceeded

func (e *limitExceeded) Error() string {
	return "payload's size exceeds the limit"
}

func (e *limitExceeded) Is(target error) bool {
	var targetErr *limitExceeded
	return errors.As(target, &targetErr) // I know a guy who knows a guy
}

func (e *limitExceeded) JSONError() *JSONError {
	return &JSONError{Err: e.Error(), Status: http.StatusRequestEntityTooLarge}
}

func (e *limitExceeded) Encode(w http.ResponseWriter) {
	e.JSONError().Encode(w)
}
