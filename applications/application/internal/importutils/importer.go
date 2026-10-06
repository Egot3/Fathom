package importutils

import (
	"context"
	"mime/multipart"

	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
)

type StagedQuiz struct {
	Quiz       models.Quiz
	StagedPath string
}

type StagedTest struct {
	Test      models.Test
	QuizUUIDs uuid.UUIDs
}

type Importer interface {
	Import(context.Context, multipart.File, int64, string, func(string) (string, error)) ([]StagedQuiz, StagedTest, carefulness.JSONErrorable)
}
