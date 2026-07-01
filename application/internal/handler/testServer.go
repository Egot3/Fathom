package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/contracts"
	"github.com/egot3/fathom/internal/logging"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// AddQuizzes implements [Service].
func (c *chiService) AddQuizzes(w http.ResponseWriter, r *http.Request) {
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

	var req contracts.AddQuizzesToTestRequest
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

	err = c.testRepo.BundleQuizzesToTest(ctx, testUUID, req.QuizUUIDs)
	if err != nil {
		logger.Info("couldn't extend test",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "couldn't create test"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *chiService) AddQuizzesToRunning(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)
	w.Header().Set("Content-Type", "application/json")

	var req contracts.AddQuizzesToTestRequest
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

	dedupUUIDs := lo.FindDuplicates(req.QuizUUIDs)

	pathes, err := c.testRepo.TestPathes(ctx, dedupUUIDs)
	if err != nil {
		logger.Info("get pathes for quizzes",
			slog.String("Error", err.Error()),
		)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "couldn't find all pathes for quizzes"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "couldn't get pathes for quizzes"})
		return
	}

	err = c.runner.UpsertQuiz(pathes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)                                 // all returned errors are user dependant anyways
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: err.Error()}) // all errors are user readable anyway
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteTest implements [Service].
func (c *chiService) DeleteTest(w http.ResponseWriter, r *http.Request) {
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

	err := c.testRepo.DeleteTest(ctx, testUUID)
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

	w.WriteHeader(http.StatusNoContent)
}

// ExtendTest implements [Service].
func (c *chiService) ExtendTest(w http.ResponseWriter, r *http.Request) {
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

	var req contracts.ExtendTestRequest
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
	logger = logger.With(slog.String("extend_by", req.ExtendBy))
	ctx = logging.WithLogger(ctx, logger)

	extendBy, err := time.ParseDuration(req.ExtendBy)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: err.Error()}) //always parseError
		return
	}
	err = c.runner.ExtendTime(extendBy)
	if err != nil {
		w.WriteHeader(http.StatusLocked)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: err.Error()}) // my own error
	}

	w.WriteHeader(http.StatusNoContent)
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

	var req contracts.PatchTestRequest
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
	logger = logger.With(slog.String("test_uuid", testUUID.String()))
	ctx = logging.WithLogger(ctx, logger)

	if req.Name == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	name := *req.Name
	if len(name) < 3 {
		logger.Info("Attempt to create testt with invalid nickname",
			slog.String("test_name", name),
		)
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "test name is too short"})
		return
	}
	if len(name) > 255 {
		logger.Info("Attempt to create test with invalid nickname",
			slog.String("test_name", name),
		)
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "test name is too big"})
		return
	}

	err = c.testRepo.UpdateTest(ctx, testUUID, name)
	if err != nil {
		logger.Info("couldn't create test",
			slog.String("Error", err.Error()),
		)
		if conflict, ok := errors.AsType[carefulness.Conflict](err); ok {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(conflict.JSONError())
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "requested test not found"})
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "couldn't create test"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PauseTest implements [Service].
func (c *chiService) PauseTest(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	// there is like nothing to return
	w.Header().Set("Content-Type", "application/json")

	err := c.runner.Pause()
	if err != nil {
		logger.Error("couldn't pause test", slog.String("Error", err.Error()))
		w.WriteHeader(http.StatusLocked)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

func (c *chiService) RemoveQuizzesFromRunning(w http.ResponseWriter, r *http.Request) {
	panic("guess it")
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
