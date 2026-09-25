package carefulness

import (
	"encoding/json"
	"errors"
	"net/http"
)

type UnprocessableRequest struct {
}

var ErrUnprocessableRequest UnprocessableRequest

func (e UnprocessableRequest) Error() string {
	return "Provided JSON is badly typed"
}

func (e UnprocessableRequest) JSONError() JSONError {
	return JSONError{Err: e.Error(), Status: http.StatusUnprocessableEntity}
}

func (e UnprocessableRequest) Is(target error) bool {
	var err *json.UnsupportedTypeError
	return errors.As(target, &err)
}

func (e UnprocessableRequest) Encode(w http.ResponseWriter) {
	e.JSONError().Encode(w)
}
