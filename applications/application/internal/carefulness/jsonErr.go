package carefulness

import (
	"encoding/json"
	"net/http"
)

type JSONError struct {
	Err    string `json:"error"`
	Status int
}

type JSONErrorable interface {
	error
	JSONError() JSONError
	Encode(http.ResponseWriter)
}

func (e JSONError) Encode(w http.ResponseWriter) {
	w.WriteHeader(e.Status)
	json.NewEncoder(w).Encode(e)
}

func (e JSONError) Error() string {
	return e.Err
}

func (e JSONError) JSONError() JSONError {
	return e // this looks whicked
}
