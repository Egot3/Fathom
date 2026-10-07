package importutils

import (
	"context"
	"mime/multipart"
	"path"
	"path/filepath"

	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
)

type StagedQuiz struct {
	Quiz       models.Quiz
	StagedPath string
	RelPath    string
}

func cleanName(n string) string { return path.Clean(filepath.ToSlash(n)) }

type StagedTest struct {
	Test      models.Test
	QuizUUIDs uuid.UUIDs
	QuizPaths []string
}

type Importer interface {
	Import(context.Context, multipart.File, int64, string, func(string) (string, error)) ([]StagedQuiz, StagedTest, carefulness.JSONErrorable)
}
