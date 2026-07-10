package carefulness

import (
	"errors"
	"fmt"
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
	return JSONError{Error: e.Error()}
}
