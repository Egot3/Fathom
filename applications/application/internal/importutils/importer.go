package importutils

import (
	"context"
	"mime/multipart"

	"github.com/egot3/fathom/internal/models"
)

type StagedQuiz struct {
	Quiz       models.Quiz
	StagedPath string
}

type Importer interface {
	Import(ctx context.Context, r multipart.File, size int64, stageDir string) ([]StagedQuiz, error)
}
