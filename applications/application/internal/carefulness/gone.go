package carefulness

import (
	"net/http"
)

type Gone struct {
}

var ErrGone Gone

func (e Gone) Error() string {
	return "Gone"
}

func (e Gone) JSONError() JSONError {
	return JSONError{Err: e.Error(), Status: http.StatusGone}
}

func (e Gone) Encode(w http.ResponseWriter) {
	e.JSONError().Encode(w)
}
