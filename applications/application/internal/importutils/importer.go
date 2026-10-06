package importutils

import (
	"context"
	"mime/multipart"

	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/models"
)

type StagedQuiz struct {
	Quiz       models.Quiz
	StagedPath string
}

type Importer interface {
	Import(context.Context, multipart.File, int64, string, func(string) (string, error)) ([]StagedQuiz, carefulness.JSONErrorable)
}
