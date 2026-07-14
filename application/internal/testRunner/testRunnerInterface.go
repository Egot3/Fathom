package testrunner

import (
	"context"
	"time"

	"github.com/egot3/fathom/internal/quiz"
	"github.com/google/uuid"
)

type TestRunner interface {
	Start(ctx context.Context, duration time.Duration, quizPaths []string, quizUUIDs, groupUUIDs uuid.UUIDs, testUUID uuid.UUID) error
	Get(quizUUID uuid.UUID) (*quiz.Quiz, error)
	cleanup(gen uint64)
	Stop()
	UpsertQuiz(quizPaths []string, quizUUIDs uuid.UUIDs) error
	RemoveQuiz(UUIDs uuid.UUIDs) error
	ExtendTime(duration time.Duration) error
	Resume() error
	Pause() error
	CurrentTestUUID() uuid.UUID
	Deadline() (*time.Time, error)

	AllowedGroupUUIDs() uuid.UUIDs
	GetAll() (uuid.UUIDs, uint64)
}
