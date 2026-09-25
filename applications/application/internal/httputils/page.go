package httputils

import (
	"net/http"
	"strconv"

	"github.com/egot3/fathom/internal/carefulness"
)

func Page(r *http.Request) (int, int, carefulness.JSONErrorable) {
	err := r.ParseForm()
	if err != nil {
		return 0, 0, carefulness.JSONError{Err: "Failed to parse form data", Status: http.StatusBadRequest}
	}

	pageInt, err := strconv.Atoi(r.Form.Get("page"))
	if err != nil {
		return 0, 0, carefulness.JSONError{Err: "given form page is not a number", Status: http.StatusUnprocessableEntity}
	}
	if pageInt < 0 {
		return 0, 0, carefulness.JSONError{Err: "page can't be < 0", Status: http.StatusUnprocessableEntity}
	}

	sizeInt, err := strconv.Atoi(r.Form.Get("size"))
	if err != nil {
		return 0, 0, carefulness.JSONError{Err: "given form size is not a number", Status: http.StatusUnprocessableEntity}
	}
	if sizeInt <= 0 {
		return 0, 0, carefulness.JSONError{Err: "page can't be <= 0", Status: http.StatusUnprocessableEntity}
	}

	return pageInt, sizeInt, nil
}
