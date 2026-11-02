package carefulness

import (
	"net/http"
)

type MalformedRequest struct {
}

var ErrMalformedRequest *MalformedRequest

func (e *MalformedRequest) Error() string {
	return "Provided JSON is malformed"
}

func (e *MalformedRequest) JSONError() *JSONError {
	return &JSONError{Err: e.Error(), Status: http.StatusBadRequest}
}

func (e *MalformedRequest) Encode(w http.ResponseWriter) {
	e.JSONError().Encode(w)
}
