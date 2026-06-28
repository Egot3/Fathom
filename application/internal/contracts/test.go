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
