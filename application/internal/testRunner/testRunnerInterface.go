package testrunner

import (
	"context"
	"time"

	"github.com/egot3/fathom/internal/quiz"
)

type TestRunner interface {
	Start(ctx context.Context, duration time.Duration, quizPaths []string) error
	Get(id int) (*quiz.Quiz, error)
	cleanup(gen uint64)
	Stop()
	UpsertQuiz(quizPaths []string) error
	RemoveQuiz(ids []int) error
	ExtendTime(duration time.Duration) error
	Resume() error
	Pause() error
}
