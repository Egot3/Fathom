package answer

import (
	"context"

	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
)

type AnswerRepository interface {
	SetAnswer(ctx context.Context, answer models.Answer) error
	Totalize(ctx context.Context, userUUID, testUUID, groupUUID uuid.UUID) error
	AnswerScore(ctx context.Context, userUUID, testUUID, groupUUID uuid.UUID, quizPath string) (int, error)
	Totals(ctx context.Context, userUUID, testUUID uuid.UUID) ([]models.UserGroupsTests, error)
}
