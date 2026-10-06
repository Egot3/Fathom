package importutils

import (
	"archive/zip"
	"context"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/egot3/fathom/internal/carefulness"
	exportutils "github.com/egot3/fathom/internal/exportUtils"
	"github.com/egot3/fathom/internal/httputils"
	"github.com/egot3/fathom/internal/logging"
	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.yaml.in/yaml/v4"
)

type zipImporter struct{}

func NewZipImporter() Importer {
	return zipImporter{}
}

func (z zipImporter) Import(ctx context.Context, r multipart.File, size int64, stageDir string, turnToAbs func(string) (string, error)) ([]StagedQuiz, StagedTest, carefulness.JSONErrorable) {
	logger := logging.LoggerFromContext(ctx).With(slog.String("strategy", httputils.Zip))

	zipReader, err := zip.NewReader(r, size)
	if err != nil {
		logger.Error("couldn't create new zip-reader", slog.String("Error", err.Error()))
		return nil, StagedTest{}, carefulness.JSONError{Err: "failed to start reading zip", Status: http.StatusUnprocessableEntity}
	}

	staged := make([]StagedQuiz, len(zipReader.File))
	stagedTest := StagedTest{}
	for i, f := range zipReader.File {
		relPath := f.Name
		stagePath := filepath.Join(stageDir, relPath)
		absPath, err := turnToAbs(relPath)
		if err != nil {
			logger.Error("couldn't get abs path", slog.String("Error", err.Error()))
			return nil, StagedTest{}, carefulness.JSONError{Err: "couldn't define writing path", Status: http.StatusInternalServerError}
		}
		if !strings.HasPrefix(absPath, filepath.Clean(stageDir)+string(os.PathSeparator)) {
			logger.Error("real zip-slip")
			return nil, StagedTest{}, carefulness.ErrZipSlip
		}

		if f.FileInfo().IsDir() {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			logger.Error("couldn't create reader for file from zip reader", slog.String("Error", err.Error()))
			return nil, StagedTest{}, carefulness.JSONError{Err: "can't read individual file", Status: http.StatusUnprocessableEntity}
		}
		defer rc.Close()

		switch filepath.Ext(f.Name) {
		case ".md":
			q, jerr := writeValidQuiz(ctx, rc, stagePath, absPath)
			if jerr != nil {
				logger.Error("couldn't write quiz", slog.String("Error", err.Error()))
				return nil, StagedTest{}, jerr
			}
			staged[i] = q
		case ".yaml":
			var t exportutils.YamlTest
			err := yaml.NewDecoder(rc).Decode(&t)
			if err != nil {
				logger.Error("couldn't unmarshal potential test", slog.String("Error", err.Error()))
				return nil, StagedTest{}, carefulness.JSONError{Err: "badly formatted test", Status: http.StatusUnprocessableEntity}
			}
			stagedTest.Test = models.Test{
				UUID: t.UUID,
				Name: t.Name,
			}
			stagedTest.QuizUUIDs = lo.Map(t.Quizzes, func(q exportutils.YamlQuiz, _ int) uuid.UUID {
				return q.UUID
			})
		}
	}

	return staged, stagedTest, nil
}
