package handler

import (
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/config"
	"github.com/egot3/fathom/internal/contracts"
	"github.com/egot3/fathom/internal/logging"
	quizparser "github.com/egot3/fathom/internal/quizParser"
	"github.com/google/uuid"
	"github.com/zeebo/xxh3"
	"go.yaml.in/yaml/v4"
)

// DeleteQuiz implements [Service].
func (c *chiService) DeleteQuiz(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")

	quizUUID, ok := (r.Context().Value("uuid")).(uuid.UUID)
	if !ok {
		logger.Error("Bad uuid")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Unable to retrieve uuid"})
		return
	}

	ctx = logging.WithLogger(ctx, logger.With(
		slog.String("quizUUID", quizUUID.String()),
	))

	err := c.quizRepo.DeallocateQuiz(ctx, quizUUID)
	if err != nil {
		logger.Error("unable to deallocate quiz",
			slog.String("Error", err.Error()),
		)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Quiz not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if s := r.Header.Get("force-delete"); s == "true" {
		quizPath, err := c.quizRepo.QuizPath(ctx, quizUUID)
		if err != nil {
			logger.Error("unable to retrieve quizPath",
				slog.String("Error", err.Error()),
			)
			w.WriteHeader(http.StatusMultiStatus)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "unable to get test's path internally"})
			return
		}

		err = os.Remove(quizPath)
		if err != nil {
			logger.Error("unable to delete quiz by quizPath",
				slog.String("Error", err.Error()),
			)
			w.WriteHeader(http.StatusMultiStatus)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "unable to delete quiz"})
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetQuiz implements [Service].
// teacher-only path
func (c *chiService) GetQuiz(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")

	quizUUID, ok := (r.Context().Value("uuid")).(uuid.UUID)
	if !ok {
		logger.Error("Bad uuid")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Unable to retrieve uuid"})
		return
	}

	ctx = logging.WithLogger(ctx, logger.With(
		slog.String("quizUUID", quizUUID.String()),
	))

	quizPath, err := c.quizRepo.QuizPath(ctx, quizUUID)
	if err != nil {
		logger.Error("unable to retrieve quizPath",
			slog.String("Error", err.Error()),
		)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Quiz not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buf, err := os.ReadFile(quizPath)
	if err != nil {
		logger.Error("Failed to read file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	quiz, err := quizparser.ParseQuizByBytes(buf)
	if err != nil {
		logger.Error("Unable to retrieve quiz by path",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contracts.GetQuizResponse{Quiz: *quiz})
}

// ListQuiz implements [Service].
func (c *chiService) ListQuizzes(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostQuiz implements [Service].
func (c *chiService) PostQuiz(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")
	var req contracts.PostQuizRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error("error in register during reading",
			slog.String("error", err.Error()),
		)
		if errors.Is(err, carefulness.ErrMalformedRequest) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.ErrMalformedRequest.JSONError()) // no err check as рукописи не горят

			return
		}
		if errors.Is(err, carefulness.ErrUnprocessableRequest) {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(carefulness.ErrUnprocessableRequest.JSONError())

			return
		}
		if errors.Is(err, io.EOF) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Empty body"})

			return
		}
		if errors.Is(err, io.ErrUnexpectedEOF) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Data loss"})

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if req.Name == "" || req.Body == "" {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error("request contains no quiz's body or its name",
			slog.String("name", req.Name),
			slog.Int("body len", len(req.Body)),
		)

		return
	}

	logger = logger.With(slog.String("name", req.Name))
	ctx = logging.WithLogger(ctx, logger)

	frontmatter, err := yaml.Marshal(req.Meta)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error("unable to process frontmatter",
			slog.String("Error", err.Error()),
		)

		return
	}
	buf := []byte(fmt.Sprintf(`
	---
	%v
	---
	%v
	`, string(frontmatter), req.Body))

	_, err = quizparser.ParseQuizByBytes(buf)
	if err != nil {
		logger.Error("couldn't parse quiz", slog.String("Error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: err.Error()})
		return
	}

	abs := config.TurnToAbs(req.Name)
	checksumUint := xxh3.HashString(req.Body)
	checksum := binary.BigEndian.AppendUint64(nil, checksumUint)
	err = c.quizRepo.RegisterQuiz(ctx, abs, checksum, req.Meta.Score)
	if err != nil {
		logger.Error("couldn't register quiz", slog.String("Error", err.Error()))
		if conflict, ok := errors.AsType[carefulness.Conflict](err); ok {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(conflict.JSONError())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "couldn't register quiz"})
		return
	}

	err = os.WriteFile(abs, buf, os.ModeAppend)
	if err != nil {
		w.WriteHeader(http.StatusMultiStatus)

		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "unable to write file"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PatchQuiz implements [Service].
func (c *chiService) PatchQuiz(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
