package answer

import (
	"context"

	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
)

type AnswerRepository interface {
	SetAnswer(ctx context.Context, testUUID, groupUUID, userUUID uuid.UUID, quizPath, answerValue string, score int) error
	Totalize(ctx context.Context, userUUID, testUUID, groupUUID uuid.UUID) error
	AnswerScore(ctx context.Context, userUUID, testUUID, groupUUID uuid.UUID, quizPath string) (int, error)
	GroupTestTotals(ctx context.Context, testUUID, groupUUID uuid.UUID) ([]models.UserGroupsTests, error)
	TestTotals(ctx context.Context, testUUID uuid.UUID) ([]models.UserGroupsTests, error)
	AllTotals(ctx context.Context, userUUID uuid.UUID) ([]models.UserGroupsTests, error)
	Total(ctx context.Context, userUUID, testUUID, groupUUID uuid.UUID) ([]models.UserGroupsTests, error)
	Answer(ctx context.Context, userUUID, testUUID, groupUUID uuid.UUID, quizPath string) (string, error)
}
