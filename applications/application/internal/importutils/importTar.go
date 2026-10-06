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
	"go.yaml.in/yaml/v4"
)

func importTar(tarReader *tar.Reader, ctx context.Context, stageDir string, turnToAbs func(string) (string, error)) ([]StagedQuiz, StagedTest, carefulness.JSONErrorable) {
	logger := logging.LoggerFromContext(ctx)

	var (
		quizzes  []StagedQuiz
		manifest *exportutils.YamlTest
	)
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
		absPath, err := turnToAbs(strings.TrimSuffix(relPath, ".md"))
		if err != nil {
			logger.Error("couldn't get abs path", slog.String("Error", err.Error()))
			return nil, StagedTest{}, carefulness.JSONError{Err: "couldn't define writing path", Status: http.StatusInternalServerError}
		}
		if !strings.HasPrefix(stagePath, filepath.Clean(stageDir)+string(os.PathSeparator)) {
			logger.Error("real gzip-slip")
			return nil, StagedTest{}, carefulness.ErrZipSlip
		}

		if f.FileInfo().IsDir() {
			continue
		}
		rel := cleanName(f.Name)

		switch filepath.Ext(rel) {
		case ".md":
			single, jerr := writeValidQuiz(ctx, io.NopCloser(tarReader), stagePath, absPath)
			if jerr != nil {
				return nil, StagedTest{}, jerr
			}
			single.RelPath = rel
			quizzes = append(quizzes, single)
		case ".yaml":
			if manifest != nil {
				return nil, StagedTest{}, carefulness.JSONError{Err: "multiple manifests in archive", Status: http.StatusUnprocessableEntity}
			}
			var t exportutils.YamlTest

			if err := yaml.NewDecoder(tarReader).Decode(&t); err != nil {
				logger.Error("couldn't unmarshal potential test", slog.String("Error", err.Error()))
				return nil, StagedTest{}, carefulness.JSONError{Err: "badly formatted test", Status: http.StatusUnprocessableEntity}
			}
			manifest = &t
		}
	}

	return reconcile(quizzes, manifest)
}
