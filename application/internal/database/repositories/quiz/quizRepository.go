package quiz

import (
	"context"

	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
)

type QuizRepository interface {
	RegisterQuiz(ctx context.Context, path string, checksum []byte, score int) error
	DeallocateQuiz(ctx context.Context, quizUUID uuid.UUID) error
	ListQuizzes(ctx context.Context, page, size int) ([]models.Quiz, int, error)
	CheckRegistered(ctx context.Context, path string) (bool, error)
	CheckIntegrity(ctx context.Context, path string, checksum []byte) (bool, error)
	QuizPath(ctx context.Context, quizUUID uuid.UUID) (string, error)
	// fuzzy search goes to app layer, as sqlite is... for storage only
	UpdateChecksum(ctx context.Context, quizUUID uuid.UUID, checksum []byte) error
	PatchQuiz(ctx context.Context, quizUUID uuid.UUID, path *string, score *int) error
}
