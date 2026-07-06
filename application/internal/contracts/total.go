package contracts

import (
	"time"

	"github.com/egot3/fathom/internal/models"
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

type GetAnswerResponse struct {
	Answer Answer `json:"answer"`
}

type GetAnswersResponse struct {
	Answers []Answer `json:"answers"`
}

type PostAnswerRequest struct {
	Value quiz.QuizAnswers `json:"answer"`
}

type Totals struct {
	Totals []models.UserGroupsTests `json:"totals"`
}

type Total struct {
	Total models.UserGroupsTests `json:"total"`
}
