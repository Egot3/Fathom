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
	GroupUUID   uuid.UUID `json:"group_uuid" bun:"group_uuid"`
	TestUUID    uuid.UUID `json:"test_uuid" bun:"test_uuid"`
	UserUUID    uuid.UUID `json:"user_uuid" bun:"user_uuid"`
	QuizUUID    uuid.UUID `json:"quiz_uuid" bun:"quiz_uuid"`
	Chosen      string    `json:"chosen" bun:"answer_value"`
	Correct     string    `json:"correct" bun:"correct"`
	SubmittedAt time.Time `json:"submitted_at" bun:"answered_at"`

	GroupName string `json:"group_name" bun:"group_name"`
	TestName  string `json:"test_name" bun:"test_name"`
}

type Total struct {
	GroupUUID uuid.UUID `json:"group_uuid" bun:"group_uuid"`
	TestUUID  uuid.UUID `json:"test_uuid" bun:"test_uuid"`
	UserUUID  uuid.UUID `json:"user_uuid" bun:"user_uuid"`
	GroupName string    `json:"group_name" bun:"group_name"`
	TestName  string    `json:"test_name" bun:"test_name"`
	Score     float64   `json:"score" bun:"score"`
}

type AnswerResponse struct {
	Answer Answer `json:"answer"`
}

type AnswersResponse struct {
	Answers []Answer `json:"answers"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	Size    int      `json:"size"`
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
