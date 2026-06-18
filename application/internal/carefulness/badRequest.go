package carefulness

import (
	"encoding/json"
	"errors"
)

type MalformedRequest struct {
}

var ErrMalformedRequest MalformedRequest

func (e MalformedRequest) Error() string {
	return "Provided JSON is malformed"
}

func (e MalformedRequest) JSONError() JSONError {
	return JSONError{Error: e.Error()}
}

func (e MalformedRequest) Is(target error) bool {
	var err *json.SyntaxError
	return errors.As(target, &err)
}
