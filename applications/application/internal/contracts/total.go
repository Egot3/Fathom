package contracts

import (
	"time"

	"github.com/egot3/fathom/internal/quiz"
	"github.com/google/uuid"
)

/*
	user, group, test go as:
	- /group/user/test/quiz/
*/

type Answer struct {
	GroupUUID   uuid.UUID `json:"group_uuid"`
	TestUUID    uuid.UUID `json:"test_uuid"`
	UserUUID    uuid.UUID `json:"user_uuid"`
	QuizUUID    uuid.UUID `json:"quiz_uuid"`
	Chosen      string    `json:"chosen"`
	Correct     string    `json:"correct"`
	SubmittedAt time.Time `json:"submitted_at"`
}

type Total struct {
	GroupUUID   uuid.UUID `json:"group_uuid" bun:"group_uuid"`
	TestUUID    uuid.UUID `json:"test_uuid" bun:"test_uuid"`
	UserUUID    uuid.UUID `json:"user_uuid" bun:"user_uuid"`
	GroupName   string    `json:"group_name" bun:"group_name"`
	TestName    string    `json:"test_name" bun:"test_name"`
	FinalizedAt time.Time `json:"finalized_at" bun:"finalized_at"`
	Score       float64   `json:"score" bun:"score"`
}

type GetAnswerResponse struct {
	Answer Answer `json:"answer"`
}

type GetAnswersResponse struct {
	Answers []Answer `json:"answers"`
}

type PostAnswerRequest struct {
	Value quiz.QuizAnswers `json:"answer"`
}

type TotalsResponse struct {
	Totals []Total `json:"totals"`
	Total  int     `json:"total"`
	Page   int     `json:"page"`
	Size   int     `json:"size"`
}

type TotalResponse struct {
	Total Total `json:"total"`
}
