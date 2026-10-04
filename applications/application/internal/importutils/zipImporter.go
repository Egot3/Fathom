package importutils

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/egot3/fathom/internal/httputils"
	"github.com/egot3/fathom/internal/logging"
	"github.com/egot3/fathom/internal/models"
	quizparser "github.com/egot3/fathom/internal/quizParser"
	"github.com/zeebo/xxh3"
)

type zipImporter struct{}

func NewZipImporter() Importer {
	return zipImporter{}
}

func (z zipImporter) Import(ctx context.Context, r multipart.File, size int64, stageDir string) ([]StagedQuiz, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.String("strategy", httputils.Zip))

	zipReader, err := zip.NewReader(r, size)
	if err != nil {
		logger.Error("couldn't create new zip-reader", slog.String("Error", err.Error()))
		return nil, err
	}

	staged := make([]StagedQuiz, len(zipReader.File))
	for i, f := range zipReader.File {
		cleanPath := filepath.Clean(f.FileInfo().Name())
		destPath := filepath.Join(stageDir, cleanPath)
		absPath, err := filepath.Abs(destPath)
		if err != nil {
			logger.Error("zip-slip detected", slog.String("Error", err.Error()))
			return nil, err
		}
		if !strings.HasPrefix(absPath, filepath.Clean(stageDir)+string(os.PathSeparator)) {
			logger.Error("real zip-slip")
			return nil, err
		}

		if f.FileInfo().IsDir() {
			if err := os.Mkdir(absPath, 0o750); err != nil {
				logger.Error("couldn't create dir", slog.String("Error", err.Error()))
				if errors.Is(err, os.ErrExist) {
					continue
				}
				return nil, err
			}
			continue
		}

		rc, err := f.Open()
		if err != nil {
			logger.Error("couldn't create reader for file from zip reader", slog.String("Error", err.Error()))
			return nil, err
		}
		defer rc.Close()

		var buf bytes.Buffer

		tee := io.TeeReader(rc, &buf)

		q, err := quizparser.ParseQuiz(tee)
		if err != nil {
			logger.Error("invalid quiz", slog.String("Error", err.Error()))
			return nil, err
		}

		dest, err := os.Create(absPath)
		if err != nil {
			logger.Error("couldn't create file for zip file", slog.String("Error", err.Error()))
			return nil, err
		}
		_, err = io.Copy(dest, &buf)
		if err != nil {
			logger.Error("couldn't write zip entry to file", slog.String("Error", err.Error()))
			return nil, err
		}
		rc.Close()
		dest.Close()

		checksumUint := xxh3.Hash(buf.Bytes())
		checksum := [8]byte(binary.BigEndian.AppendUint64(nil, checksumUint))

		answer, err := json.Marshal(q.Answer)
		if err != nil {
			logger.Error("couldn't marshal answer to json",
				slog.String("Error", err.Error()),
			)
			return nil, err
		}

		staged[i] = StagedQuiz{
			StagedPath: destPath,
			Quiz: models.Quiz{
				UUID:          q.UUID, // wrong, placeholder
				Path:          absPath,
				Checksum:      checksum,
				CorrectAnswer: string(answer),
				Score:         q.Meta.Score,
			},
		}
	}

	return staged, nil
}
