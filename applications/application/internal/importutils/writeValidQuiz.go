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
	"path/filepath"

	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/logging"
	"github.com/egot3/fathom/internal/models"
	quizparser "github.com/egot3/fathom/internal/quizParser"
	"github.com/zeebo/xxh3"
)

func writeValidQuiz(ctx context.Context, reader io.ReadCloser, stagePath, absPath string) (StagedQuiz, carefulness.JSONErrorable) {
	logger := logging.LoggerFromContext(ctx).With(slog.String("type", "quiz"))

	raw, err := io.ReadAll(io.LimitReader(reader, 1<<20))
	q, err := quizparser.ParseQuiz(bytes.NewReader(raw))
	checksumUint := xxh3.Hash(raw)
	_ = os.MkdirAll(filepath.Dir(stagePath), 0o755)
	err = os.WriteFile(stagePath, raw, 0o644)

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
