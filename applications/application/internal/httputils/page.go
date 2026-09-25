package httputils

import (
	"errors"
	"net/http"
	"strconv"
)

func Page(r *http.Request) (int, int, error) {
	err := r.ParseForm()
	if err != nil {
		return 0, 0, errors.New("Failed to parse form data")
	}

	pageInt, err := strconv.Atoi(r.Form.Get("page"))
	if err != nil {
		return 0, 0, errors.New("given form page is not a number")
	}
	if pageInt < 0 {
		return 0, 0, errors.New("page can't be < 0")
	}

	sizeInt, err := strconv.Atoi(r.Form.Get("size"))
	if err != nil {
		return 0, 0, errors.New("given form size is not a number")
	}
	if sizeInt <= 0 {
		return 0, 0, errors.New("size can't be <= 0")
	}

	return pageInt, sizeInt, nil
}
