package contracts

import (
	"github.com/egot3/fathom/internal/models"
	"github.com/egot3/fathom/internal/quiz"
	"github.com/google/uuid"
)

type GetQuizResponse struct {
	Meta quiz.Frontmatter `json:"meta"`
	Body string           `json:"body"`
}

type PostQuizRequest struct {
	Name string           `json:"name"`
	Body string           `json:"body"`
	Meta quiz.Frontmatter `json:"meta"`
}

type PatchQuizRequest struct {
	Name *string           `json:"name"`
	Body *string           `json:"body"`
	Meta *quiz.Frontmatter `json:"meta"`
}

type ListQuizResponse struct {
	Quizzes []models.Quiz `json:"quizzes"`
	Total   int           `json:"total"`
	Page    int           `json:"page"`
	Size    int           `json:"size"`
}

type ExportQuizRequest struct {
	UUIDs uuid.UUIDs `json:"uuids"`
}

type ParsedQuizResponse struct {
	Quiz quiz.Quiz `json:"quiz"`
}
