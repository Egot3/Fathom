package contracts

import "github.com/egot3/fathom/internal/quiz"

type GetQuizResponse struct {
	Quiz quiz.Quiz `json:"quiz"`
}

type PostQuizRequest struct {
	Name string           `json:"name"`
	Body string           `json:"body"`
	Meta quiz.Frontmatter `json:"meta"`
}
