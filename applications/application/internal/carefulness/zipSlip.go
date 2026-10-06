package carefulness

import (
	"errors"
	"net/http"
)

type ZipSlip struct{}

var ErrZipSlip ZipSlip

func (e ZipSlip) Error() string {
	return "Zip slip in archive detected"
}

func (e ZipSlip) JSONError() JSONError {
	return JSONError{Err: e.Error(), Status: http.StatusBadRequest}
}

func (e ZipSlip) Is(target error) bool {
	return errors.As(target, &e)
}

func (e ZipSlip) Encode(w http.ResponseWriter) {
	e.JSONError().Encode(w)
}
