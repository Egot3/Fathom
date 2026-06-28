package contracts

import (
	"github.com/egot3/fathom/internal/models"
	"github.com/egot3/fathom/internal/quiz"
)

type GetQuizResponse struct {
	Quiz quiz.Quiz `json:"quiz"`
}

type PostQuizRequest struct {
	Name string           `json:"name"`
	Body string           `json:"body"`
	Meta quiz.Frontmatter `json:"meta"`
}

type PutQuizRequest struct {
	Body string           `json:"body"`
	Meta quiz.Frontmatter `json:"meta"`
}

type PatchQuizRequest struct {
	Name  *string `json:"name"`
	Score *int    `json:"score"`
}

type ListQuizResponse struct {
	Quizzes []models.Quiz `json:"quizzes"`
	Total   int           `json:"total"`
	Page    int           `json:"page"`
	Size    int           `json:"size"`
}
