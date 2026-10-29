package carefulness

import (
	"net/http"
)

type Panic struct{}

var ErrPanic *Panic

func (e *Panic) Error() string {
	return "an error so omnious, that even program isn't sure about what happened. If it was encountered naturally, please make an issue on our github page"
}

func (e *Panic) JSONError() *JSONError {
	return &JSONError{Err: e.Error(), Status: http.StatusInternalServerError}
}

func (e *Panic) Encode(w http.ResponseWriter) {
	e.JSONError().Encode(w)
}
