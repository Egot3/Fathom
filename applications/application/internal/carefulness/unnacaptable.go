package carefulness

import "net/http"

type Unnacaptable struct {
}

var ErrUnnacaptable Unnacaptable

func (e Unnacaptable) Error() string {
	return "no supported type found"
}

func (e Unnacaptable) JSONError() JSONError {
	return JSONError{Err: e.Error(), Status: http.StatusNotAcceptable}
}

func (e Unnacaptable) Encode(w http.ResponseWriter) {
	e.JSONError().Encode(w)
}
