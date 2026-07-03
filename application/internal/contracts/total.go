package contracts

import (
	"time"

	"github.com/egot3/fathom/internal/quiz"
	"github.com/google/uuid"
)

/*
	answer, user, group, test go as:
	- /group/user/test/quiz/answer
*/

type Answer struct {
	GroupUUID   uuid.UUID        `json:"group_uuid"`
	TestUUID    uuid.UUID        `json:"test_uuid"`
	UserUUID    uuid.UUID        `json:"user_uuid"`
	QuizUUID    uuid.UUID        `json:"quiz_uuid"`
	Chosen      quiz.QuizOptions `json:"chosen"`
	Right       quiz.QuizAnswers `json:"right"`
	SubmittedAt time.Time        `json:"submitted_at"`
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
