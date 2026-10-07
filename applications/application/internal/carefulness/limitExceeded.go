package carefulness

import (
	"fmt"
	"net/http"
)

type limitExceeded struct{}

var ErrLimitExceeded limitExceeded

func (e limitExceeded) Error() string {
	return fmt.Sprintf("payload's size exceeds the limit")
}

func (e limitExceeded) Is(target error) bool {
	return target == ErrLimitExceeded
}

func (e limitExceeded) JSONError() JSONError {
	return JSONError{Err: e.Error(), Status: http.StatusRequestEntityTooLarge}
}

func (e limitExceeded) Encode(w http.ResponseWriter) {
	e.JSONError().Encode(w)
}
