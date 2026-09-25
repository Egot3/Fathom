package carefulness

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrQuizNotInTest = errors.New("quiz is not in the test")

type NotInTestError struct {
	Count int
}

func (e *NotInTestError) Error() string {
	return fmt.Sprintf("%d quizzes are not in test", e.Count)
}

func (e *NotInTestError) Is(target error) bool {
	return target == ErrQuizNotInTest
}

func (e *NotInTestError) JSONError() JSONError {
	return JSONError{Err: e.Error(), Status: http.StatusNotFound}
}

func (e NotInTestError) Encode(w http.ResponseWriter) {
	e.JSONError().Encode(w)
}
