package carefulness

import (
	"encoding/json"
	"errors"
	"net/http"
)

type MalformedRequest struct {
}

var ErrMalformedRequest MalformedRequest

func (e MalformedRequest) Error() string {
	return "Provided JSON is malformed"
}

func (e MalformedRequest) JSONError() JSONError {
	return JSONError{Err: e.Error(), Status: http.StatusBadRequest}
}

func (e MalformedRequest) Is(target error) bool {
	var err *json.SyntaxError
	return errors.As(target, &err)
}

func (e MalformedRequest) Encode(w http.ResponseWriter) {
	e.JSONError().Encode(w)
}
