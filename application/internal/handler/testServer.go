package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/contracts"
	"github.com/egot3/fathom/internal/logging"
	"github.com/google/uuid"
)

// AddQuizzes implements [Service].
func (c *chiService) AddQuizzes(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func (c *chiService) AddQuizzesToRunning(w http.ResponseWriter, r *http.Request) {
	/* logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)
	w.Header().Set("Content-Type", "application/json") */
}

// DeleteTest implements [Service].
func (c *chiService) DeleteTest(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// ExtendTest implements [Service].
func (c *chiService) ExtendTest(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// GetTest implements [Service].
func (c *chiService) GetTest(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)
	w.Header().Set("Content-Type", "application/json")

	testUUID, ok := (r.Context().Value("uuid")).(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Unable to retrieve test uuid"})
		return
	}
	logger = logger.With(slog.String("test_uuid", testUUID.String()))
	ctx = logging.WithLogger(ctx, logger)

	test, err := c.testRepo.Test(ctx, testUUID)
	if err != nil {
		logger.Error("couldn't retrive test", slog.String("Error", err.Error()))
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "requested test not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contracts.GetTestResponse{Test: *test})
}

// PatchTest implements [Service].
func (c *chiService) PatchTest(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PauseTest implements [Service].
func (c *chiService) PauseTest(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostTest implements [Service].
func (c *chiService) PostTest(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)
	w.Header().Set("Content-Type", "application/json")

	var req contracts.PostTestRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error("Failed to parse body",
			slog.String("Error", err.Error()),
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
	logger = logger.With(slog.String("test_name", req.Name))
	ctx = logging.WithLogger(ctx, logger)

	if len(req.Name) < 3 {
		logger.Info("Attempt to create testt with invalid nickname",
			slog.String("test_name", req.Name),
		)
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "test name is too short"})
		return
	}
	if len(req.Name) > 255 {
		logger.Info("Attempt to create test with invalid nickname",
			slog.String("test_name", req.Name),
		)
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "test name is too big"})
		return
	}

	test, err := c.testRepo.CreateTest(ctx, req.Name)
	if err != nil {
		logger.Info("couldn't create test",
			slog.String("Error", err.Error()),
		)
		if conflict, ok := errors.AsType[carefulness.Conflict](err); ok {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(conflict.JSONError())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "couldn't create test"})
		return
	}

	if req.Quizzes != nil {
		err := c.testRepo.BundleQuizzesToTest(ctx, test.UUID, req.Quizzes)
		if err != nil {
			logger.Info("couldn't add quizzes to test",
				slog.String("Error", err.Error()),
			)
			w.WriteHeader(http.StatusMultiStatus)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "couldn't add quizzes to test"})
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveQuizzes implements [Service].
func (c *chiService) RemoveQuizzes(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// ResumeTest implements [Service].
func (c *chiService) ResumeTest(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// StartTest implements [Service].
func (c *chiService) StartTest(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// StopTest implements [Service].
func (c *chiService) StopTest(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
