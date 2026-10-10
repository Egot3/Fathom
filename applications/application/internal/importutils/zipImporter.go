package importutils

import (
	"archive/zip"
	"context"
	"errors"
	"io"
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

	staged := make([]StagedQuiz, 0, len(zipReader.File))
	var manifest *exportutils.Manifest = nil
	for _, f := range zipReader.File {
		relPath := f.Name
		stagePath := filepath.Join(stageDir, relPath)
		absPath, err := turnToAbs(strings.TrimSuffix(relPath, ".md"))
		if err != nil {
			logger.Error("couldn't get abs path", slog.String("Error", err.Error()))
			return nil, StagedTest{}, carefulness.JSONError{Err: "couldn't define writing path", Status: http.StatusInternalServerError}
		}
		if !strings.HasPrefix(stagePath, filepath.Clean(stageDir)+string(os.PathSeparator)) {
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
				logger.Error("couldn't write quiz", slog.String("Error", jerr.Error()))
				return nil, StagedTest{}, jerr
			}
			q.RelPath = cleanName(f.Name)
			staged = append(staged, q)
		case ".yaml":
			var m exportutils.Manifest

			limited := io.LimitReader(rc, 1<<20)
			if err := yaml.NewDecoder(limited).Decode(&m); err != nil {
				logger.Error("couldn't unmarshal potential test", slog.String("Error", err.Error()))
				if errors.Is(err, io.EOF) {
					return nil, StagedTest{}, carefulness.ErrLimitExceeded
				}
				return nil, StagedTest{}, carefulness.JSONError{Err: "badly formatted test", Status: http.StatusUnprocessableEntity}
			}
			if manifest != nil {
				return nil, StagedTest{}, carefulness.JSONError{Err: "multiple manifests in archive", Status: http.StatusUnprocessableEntity}
			}

			err := m.Validate()
			if err != nil {
				return nil, StagedTest{}, carefulness.JSONError{Err: "badly formatted manifest", Status: http.StatusUnprocessableEntity}
			}

			manifest = &m
		}
	}

	return reconcile(staged, manifest, turnToAbs)
}
