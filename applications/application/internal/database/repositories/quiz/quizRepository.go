package quiz

import (
	"context"

	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
)

type QuizRepository interface {
	RegisterQuiz(ctx context.Context, path string, checksum [8]byte, score int, answer []byte) error
	DeallocateQuiz(ctx context.Context, quizUUID uuid.UUID) error
	ListQuizzes(ctx context.Context, page, size int) ([]models.Quiz, int, error)
	CheckRegistered(ctx context.Context, path string) (bool, error)
	CheckIntegrity(ctx context.Context, path string, checksum [8]byte) (bool, error)
	QuizPath(ctx context.Context, quizUUID uuid.UUID) (string, error)

	UpdateChecksum(ctx context.Context, quizUUID uuid.UUID, checksum [8]byte) error
	PatchQuiz(ctx context.Context, quizUUID uuid.UUID, path *string, score *int) error
	CorrectAnswer(ctx context.Context, quizUUID uuid.UUID) (string, error)
	ExistsByUUID(ctx context.Context, quizUUID uuid.UUID) (bool, error)
	QuizFresh(ctx context.Context, quizUUID uuid.UUID, checksum [8]byte) (bool, error)
	Quiz(ctx context.Context, quizUUID uuid.UUID) (models.Quiz, error)
}
