package httputils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/egot3/fathom/internal/carefulness"
)

// consumes r.body because, yk, it must
func ParseJSON[T any](from *http.Request, to *T) carefulness.JSONErrorable {
	if mime := from.Header.Get("Content-Type"); mime != "application/json" {
		return &carefulness.UnsupportedMediaType{Mime: mime}
	}
	err := json.NewDecoder(from.Body).Decode(to)
	if err != nil {
		if _, ok := errors.AsType[*json.SyntaxError](err); ok {
			return carefulness.ErrMalformedRequest
		}

		switch {
		case errors.Is(err, carefulness.ErrUnprocessableRequest):
			return carefulness.ErrUnprocessableRequest

		case errors.Is(err, io.EOF):
			return &carefulness.JSONError{Err: "got empty body", Status: http.StatusBadRequest}

		case errors.Is(err, io.ErrUnexpectedEOF):
			return &carefulness.JSONError{Err: "data loss", Status: http.StatusBadRequest}

		default:
			return &carefulness.JSONError{Err: "unknown error", Status: http.StatusInternalServerError}
		}
	}

	return nil
}
