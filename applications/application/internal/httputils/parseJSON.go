package httputils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/egot3/fathom/internal/carefulness"
)

func ParseJSON[T any](from io.ReadCloser, to *T) carefulness.JSONErrorable {
	err := json.NewDecoder(from).Decode(to)
	if err != nil {
		switch {
		case errors.Is(err, carefulness.ErrMalformedRequest):
			return carefulness.ErrMalformedRequest

		case errors.Is(err, carefulness.ErrUnprocessableRequest):
			return carefulness.ErrUnprocessableRequest

		case errors.Is(err, io.EOF):
			return carefulness.JSONError{Err: "got empty body", Status: http.StatusBadRequest}

		case errors.Is(err, io.ErrUnexpectedEOF):
			return carefulness.JSONError{Err: "data loss", Status: http.StatusBadRequest}

		default:
			return carefulness.JSONError{Err: "unknown error", Status: http.StatusInternalServerError}
		}
	}

	return nil
}
