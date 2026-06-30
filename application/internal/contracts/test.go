package contracts

import (
	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
)

type GetTestResponse struct {
	Test models.Test `json:"test"`
}

type PostTestRequest struct {
	Name    string     `json:"name"`
	Quizzes uuid.UUIDs `json:"quizzes"`
}

type AddQuizzesToTestRequest struct {
	QuizUUIDs uuid.UUIDs `json:"quiz_uuids"`
}

type PatchTestRequest struct {
	Name *string `json:"name"`
}

type ExtendTest struct {
	ExtendBy string `json:"extend_by"`
}
