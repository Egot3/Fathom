package carefulness

import (
	"errors"
	"fmt"
	"net/http"
)

type UnsupportedMediaType struct {
	Mime string
}

func (e *UnsupportedMediaType) Error() string {
	return fmt.Sprintf("provided media cannot be processed because of unsupported type: %q", e.Mime)
}

func (e *UnsupportedMediaType) JSONError() *JSONError {
	return &JSONError{Err: e.Error(), Status: http.StatusUnsupportedMediaType}
}

func (e *UnsupportedMediaType) Is(target error) bool {
	return errors.As(target, &e)
}

func (e *UnsupportedMediaType) Encode(w http.ResponseWriter) {
	e.JSONError().Encode(w)
}
