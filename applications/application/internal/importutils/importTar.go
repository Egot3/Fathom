package importutils

import (
	"archive/tar"
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/egot3/fathom/internal/carefulness"
	exportutils "github.com/egot3/fathom/internal/exportUtils"
	"github.com/egot3/fathom/internal/logging"
	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.yaml.in/yaml/v4"
)

func importTar(tarReader *tar.Reader, ctx context.Context, stageDir string, turnToAbs func(string) (string, error)) ([]StagedQuiz, StagedTest, carefulness.JSONErrorable) {
	logger := logging.LoggerFromContext(ctx)

	staged := []StagedQuiz{} //unknown amount of files?
	stagedTest := StagedTest{}
	for {
		f, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Error("unable to get next tar entry",
				slog.String("Error", err.Error()),
			)
			return nil, StagedTest{}, carefulness.JSONError{Err: "couldn't continue reading archive", Status: http.StatusUnprocessableEntity}
		}

		relPath := f.Name
		stagePath := filepath.Join(stageDir, relPath)
		absPath, err := turnToAbs(relPath)
		if err != nil {
			logger.Error("couldn't get abs path", slog.String("Error", err.Error()))
			return nil, StagedTest{}, carefulness.JSONError{Err: "couldn't define writing path", Status: http.StatusInternalServerError}
		}
		if !strings.HasPrefix(absPath, filepath.Clean(stageDir)+string(os.PathSeparator)) {
			logger.Error("real gzip-slip")
			return nil, StagedTest{}, carefulness.ErrZipSlip
		}

		if f.FileInfo().IsDir() {
			continue
		}

		switch filepath.Ext(f.Name) {
		case ".md":
			single, jerr := writeValidQuiz(ctx, io.NopCloser(tarReader), stagePath, absPath)
			if jerr != nil {
				return nil, StagedTest{}, jerr
			}

			staged = append(staged, single)
		case ".yaml":
			var t exportutils.YamlTest
			err := yaml.NewDecoder(tarReader).Decode(&t)
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

	return staged, StagedTest{}, nil
}
