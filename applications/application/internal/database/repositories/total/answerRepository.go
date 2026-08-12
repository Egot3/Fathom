package total

import (
	"context"

	"github.com/egot3/fathom/internal/contracts"
	"github.com/google/uuid"
)

type TotalRepository interface {
	SetAnswer(ctx context.Context, testUUID, groupUUID, userUUID, quizUUID uuid.UUID, answerValue string, score float32) error
	Totalize(ctx context.Context, userUUID, testUUID, groupUUID uuid.UUID) error
	AnswerScore(ctx context.Context, userUUID, testUUID, groupUUID, quizUUID uuid.UUID) (float32, error)
	GroupTestTotals(ctx context.Context, testUUID, groupUUID uuid.UUID) ([]contracts.Total, error)
	TestTotals(ctx context.Context, testUUID uuid.UUID) ([]contracts.Total, error)
	AllTotals(ctx context.Context, userUUID uuid.UUID, page, size int) ([]contracts.Total, int, error)
	Total(ctx context.Context, userUUID, testUUID, groupUUID uuid.UUID) (*contracts.Total, error)
	Answer(ctx context.Context, userUUID, testUUID, groupUUID, quizUUID uuid.UUID) (string, error)
	AnswersInTest(ctx context.Context, userUUID, testUUID, groupUUID uuid.UUID, page, size int) ([]contracts.Answer, int, error)
}
