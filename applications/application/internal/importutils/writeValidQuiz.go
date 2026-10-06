package importutils

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/logging"
	"github.com/egot3/fathom/internal/models"
	quizparser "github.com/egot3/fathom/internal/quizParser"
	"github.com/zeebo/xxh3"
)

func writeValidQuiz(ctx context.Context, reader io.ReadCloser, stagePath, absPath string) (StagedQuiz, carefulness.JSONErrorable) {
	logger := logging.LoggerFromContext(ctx).With(slog.String("type", "quiz"))

	var buf bytes.Buffer
	tee := io.TeeReader(reader, &buf)

	q, err := quizparser.ParseQuiz(tee)
	if err != nil {
		logger.Error("invalid quiz", slog.String("Error", err.Error()))
		return StagedQuiz{}, carefulness.JSONError{Err: "badly formatted quiz in archive", Status: http.StatusUnprocessableEntity}
	}

	dest, err := os.Create(stagePath)
	if err != nil {
		logger.Error("couldn't create file for tar file", slog.String("Error", err.Error()))
		return StagedQuiz{}, carefulness.JSONError{Err: "unable to create temp quiz", Status: http.StatusInternalServerError}
	}
	defer dest.Close()

	_, err = io.Copy(dest, &buf)
	if err != nil {
		logger.Error("couldn't write tar entry to file", slog.String("Error", err.Error()))
		return StagedQuiz{}, carefulness.JSONError{Err: "unable to write quiz to temp file", Status: http.StatusInternalServerError}
	}
	dest.Close()

	checksumUint := xxh3.Hash(buf.Bytes())
	checksum := [8]byte(binary.BigEndian.AppendUint64(nil, checksumUint))

	answer, err := json.Marshal(q.Answer)
	if err != nil {
		logger.Error("couldn't marshal answer to json",
			slog.String("Error", err.Error()),
		)
		return StagedQuiz{}, carefulness.JSONError{Err: "couldn't format correct answer", Status: http.StatusUnprocessableEntity}
	}

	return StagedQuiz{
		StagedPath: stagePath,
		Quiz: models.Quiz{
			Path:          absPath,
			Checksum:      checksum,
			CorrectAnswer: string(answer),
			Score:         q.Meta.Score,
		},
	}, nil
}
