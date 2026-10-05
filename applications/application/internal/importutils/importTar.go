package importutils

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/egot3/fathom/internal/logging"
	"github.com/egot3/fathom/internal/models"
	quizparser "github.com/egot3/fathom/internal/quizParser"
	"github.com/zeebo/xxh3"
)

func importTar(tarReader *tar.Reader, ctx context.Context, stageDir string, turnToAbs func(string) (string, error)) ([]StagedQuiz, error) {
	logger := logging.LoggerFromContext(ctx)

	staged := []StagedQuiz{} //unknown amount of files?
	for {
		f, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Error("unable to get next tar entry",
				slog.String("Error", err.Error()),
			)
			return nil, err
		}

		relPath := f.Name
		stagePath := filepath.Join(stageDir, relPath)
		absPath, err := turnToAbs(relPath)
		if err != nil {
			logger.Error("couldn't get abs path", slog.String("Error", err.Error()))
			return nil, err
		}
		if !strings.HasPrefix(absPath, filepath.Clean(stageDir)+string(os.PathSeparator)) {
			logger.Error("real gzip-slip")
			return nil, err
		}

		if f.FileInfo().IsDir() {
			continue
		}

		var buf bytes.Buffer
		tee := io.TeeReader(tarReader, &buf)

		q, err := quizparser.ParseQuiz(tee)
		if err != nil {
			logger.Error("invalid quiz", slog.String("Error", err.Error()))
			return nil, err
		}

		dest, err := os.Create(stagePath)
		if err != nil {
			logger.Error("couldn't create file for tar file", slog.String("Error", err.Error()))
			return nil, err
		}
		defer dest.Close()

		_, err = io.Copy(dest, &buf)
		if err != nil {
			logger.Error("couldn't write tar entry to file", slog.String("Error", err.Error()))
			return nil, err
		}
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

		staged = append(staged, StagedQuiz{
			StagedPath: stagePath,
			Quiz: models.Quiz{
				Path:          absPath,
				Checksum:      checksum,
				CorrectAnswer: string(answer),
				Score:         q.Meta.Score,
			},
		})

	}

	return staged, nil
}
