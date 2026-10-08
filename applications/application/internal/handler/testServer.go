package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	jwtutils "github.com/egot3/fathom/internal/JWTutils"
	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/contracts"
	exportutils "github.com/egot3/fathom/internal/exportUtils"
	"github.com/egot3/fathom/internal/httputils"
	"github.com/egot3/fathom/internal/importutils"
	"github.com/egot3/fathom/internal/logging"
	"github.com/egot3/fathom/internal/models"
	"github.com/egot3/fathom/internal/quiz"
	testrunner "github.com/egot3/fathom/internal/testRunner"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.yaml.in/yaml/v4"
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
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Unable to retrieve test uuid"})
		return
	}
	logger = logger.With(slog.String("test_uuid", testUUID.String()))
	ctx = logging.WithLogger(ctx, logger)

	var req contracts.AddQuizzesToTestRequest
	jerr := httputils.ParseJSON(r.Body, &req)
	if jerr != nil {
		logger.Error("Failed to parse body",
			slog.String("Error", jerr.Error()),
		)
		jerr.Encode(w)
		return
	}

	if len(req.QuizUUIDs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err := c.testRepo.BundleQuizzesToTest(ctx, testUUID, req.QuizUUIDs)
	if err != nil {
		logger.Error("couldn't append quizzes to test",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't append quizzes to test test"})
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
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Unable to retrieve test uuid"})
		return
	}
	logger = logger.With(slog.String("test_uuid", testUUID.String()))
	ctx = logging.WithLogger(ctx, logger)

	err := c.testRepo.DeleteTest(ctx, testUUID)
	if err != nil {
		logger.Error("couldn't retrive test", slog.String("Error", err.Error()))
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Err: "requested test not found"})
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

	runnerKey, err := strconv.ParseUint(chi.URLParam(r, "runnerKey"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.ErrMalformedRequest)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.ErrMalformedRequest)
		return
	}
	logger = logger.With(slog.Uint64("runner_key", runnerKey))
	ctx = logging.WithLogger(ctx, logger)

	var req contracts.ExtendTestRequest
	jerr := httputils.ParseJSON(r.Body, &req)
	if jerr != nil {
		logger.Error("Failed to parse body",
			slog.String("Error", jerr.Error()),
		)
		jerr.Encode(w)
		return
	}
	logger = logger.With(slog.String("extend_by", req.ExtendBy))
	ctx = logging.WithLogger(ctx, logger)

	extendBy, err := time.ParseDuration(req.ExtendBy)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: err.Error()}) //always parseError
		return
	}
	tr, ok := c.manager.Get(runnerKey)
	if !ok {
		w.WriteHeader(http.StatusLocked)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: testrunner.ErrRunnerInactive.Error()}) // my own error
		return
	}

	err = tr.ExtendTime(extendBy)
	if err != nil {
		w.WriteHeader(http.StatusLocked)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: err.Error()}) // my own error
		return
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
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Unable to retrieve test uuid"})
		return
	}
	logger = logger.With(slog.String("test_uuid", testUUID.String()))
	ctx = logging.WithLogger(ctx, logger)

	test, err := c.testRepo.Test(ctx, testUUID)
	if err != nil {
		logger.Error("couldn't retrive test", slog.String("Error", err.Error()))
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Err: "requested test not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contracts.GetTestResponse{Test: test})
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
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Unable to retrieve test uuid"})
		return
	}
	logger = logger.With(slog.String("test_uuid", testUUID.String()))
	ctx = logging.WithLogger(ctx, logger)

	var req contracts.PatchTestRequest
	jerr := httputils.ParseJSON(r.Body, &req)
	if jerr != nil {
		logger.Error("Failed to parse body",
			slog.String("Error", jerr.Error()),
		)
		jerr.Encode(w)
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "test name is too short"})
		return
	}
	if len(name) > 255 {
		logger.Info("Attempt to create test with invalid nickname",
			slog.String("test_name", name),
		)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "test name is too big"})
		return
	}

	err := c.testRepo.UpdateTest(ctx, testUUID, name)
	if err != nil {
		logger.Info("couldn't create test",
			slog.String("Error", err.Error()),
		)
		if conflict, ok := errors.AsType[carefulness.Conflict](err); ok {
			conflict.Encode(w)
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Err: "requested test not found"})
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't create test"})
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

	runnerKey, err := strconv.ParseUint(chi.URLParam(r, "runnerKey"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.ErrMalformedRequest)
		return
	}
	logger = logger.With(slog.Uint64("runner_key", runnerKey))

	tr, ok := c.manager.Get(runnerKey)
	if !ok {
		w.WriteHeader(http.StatusLocked)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: testrunner.ErrRunnerInactive.Error()})
		return
	}

	err = tr.Pause()
	if err != nil {
		logger.Error("couldn't pause test", slog.String("Error", err.Error()))
		w.WriteHeader(http.StatusLocked)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: err.Error()})
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
	jerr := httputils.ParseJSON(r.Body, &req)
	if jerr != nil {
		logger.Error("Failed to parse body",
			slog.String("Error", jerr.Error()),
		)
		jerr.Encode(w)
		return
	}
	logger = logger.With(slog.String("test_name", req.Name))
	ctx = logging.WithLogger(ctx, logger)

	if len(req.Name) < 3 {
		logger.Info("Attempt to create testt with invalid nickname",
			slog.String("test_name", req.Name),
		)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "test name is too short"})
		return
	}
	if len(req.Name) > 255 {
		logger.Info("Attempt to create test with invalid nickname",
			slog.String("test_name", req.Name),
		)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "test name is too big"})
		return
	}

	test, err := c.testRepo.CreateTest(ctx, req.Name)
	if err != nil {
		logger.Info("couldn't create test",
			slog.String("Error", err.Error()),
		)
		if conflict, ok := errors.AsType[*carefulness.Conflict](err); ok {
			conflict.Encode(w)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't create test"})
		return
	}

	if req.Quizzes != nil {
		err := c.testRepo.BundleQuizzesToTest(ctx, test.UUID, req.Quizzes)
		if err != nil {
			logger.Info("couldn't add quizzes to test",
				slog.String("Error", err.Error()),
			)
			w.WriteHeader(http.StatusMultiStatus)
			json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't add quizzes to test"})
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveQuizzes implements [Service].
func (c *chiService) RemoveQuizzes(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)
	w.Header().Set("Content-Type", "application/json")

	testUUID, ok := (r.Context().Value("uuid")).(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Unable to retrieve test uuid"})
		return
	}
	logger = logger.With(slog.String("test_uuid", testUUID.String()))
	ctx = logging.WithLogger(ctx, logger)

	var req contracts.RemoveQuizzesRequest
	jerr := httputils.ParseJSON(r.Body, &req)
	if jerr != nil {
		logger.Error("Failed to parse body",
			slog.String("Error", jerr.Error()),
		)
		jerr.Encode(w)
		return
	}

	err := c.testRepo.PruneQuizzesFromTest(ctx, testUUID, req.QuizUUIDs)
	if err != nil {
		logger.Info("couldn't prune quizzes from test",
			slog.String("Error", err.Error()),
		)
		if notFound, ok := errors.AsType[*carefulness.NotInTestError](err); ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(notFound.JSONError())
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Err: "none of the quizzes is in test"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't delete quiz from test"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ResumeTest implements [Service].
func (c *chiService) ResumeTest(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	// there is like nothing to return
	w.Header().Set("Content-Type", "application/json")

	runnerKey, err := strconv.ParseUint(chi.URLParam(r, "runnerKey"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.ErrMalformedRequest)
		return
	}
	logger = logger.With(slog.Uint64("runner_key", runnerKey))

	tr, ok := c.manager.Get(runnerKey)
	if !ok {
		w.WriteHeader(http.StatusLocked)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: testrunner.ErrRunnerInactive.Error()})
		return
	}

	err = tr.Resume()
	if err != nil {
		logger.Error("couldn't resume test", slog.String("Error", err.Error()))
		w.WriteHeader(http.StatusLocked)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// StartTest implements [Service].
func (c *chiService) StartTest(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)
	w.Header().Set("Content-Type", "application/json")

	var req contracts.StartRequest
	jerr := httputils.ParseJSON(r.Body, &req)
	if jerr != nil {
		logger.Error("Failed to parse body",
			slog.String("Error", jerr.Error()),
		)
		jerr.Encode(w)
		return
	}
	logger = logger.With(slog.String("duration", req.Duration), slog.Any("requested test", req.TestUUID))
	ctx = logging.WithLogger(ctx, logger)

	duration, err := time.ParseDuration(req.Duration)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: err.Error()}) //always parseError
		return
	}

	test, err := c.testRepo.Test(ctx, req.TestUUID)
	if err != nil {
		logger.Error("couldn't get pathes for test",
			slog.String("Error", err.Error()),
		)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't find all pathes for quizzes"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't get pathes for quizzes"})
		return
	}

	do, err := c.groupRepo.GroupsExist(ctx, req.GroupsUUIDs)
	if err != nil {
		logger.Error("couldn't check if groups exist",
			slog.String("Error", err.Error()),
		)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't check if groups exist"})
		return
	}

	if !do {
		logger.Info("some of requested groups don't exist")

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "some of requested groups don't exist"})
		return
	}

	quizPathes := make([]string, len(test.Quizzes))
	quizUUIDs := make(uuid.UUIDs, len(test.Quizzes))
	lo.ForEach(test.Quizzes, func(quiz models.Quiz, i int) {
		quizPathes[i] = quiz.Path
		quizUUIDs[i] = quiz.UUID
	})

	_, err = c.manager.Start(ctx, duration, quizPathes, quizUUIDs, req.GroupsUUIDs, test.UUID)
	if err != nil {
		logger.Error("unable to start test", slog.String("Error", err.Error()))
		w.WriteHeader(http.StatusBadRequest) // all returned errors are user dependant anyways
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// StopTest implements [Service].
func (c *chiService) StopTest(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	w.Header().Set("Content-Type", "application/json")

	runnerKey, err := strconv.ParseUint(chi.URLParam(r, "runnerKey"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.ErrMalformedRequest)
		return
	}
	logger = logger.With(slog.Uint64("runner_key", runnerKey))

	tr, ok := c.manager.Get(runnerKey)
	if !ok {
		w.WriteHeader(http.StatusLocked)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: testrunner.ErrRunnerInactive.Error()})
		return
	}

	tr.Stop()

	w.WriteHeader(http.StatusNoContent)
}

func (c *chiService) GetQuizFromRunning(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)
	w.Header().Set("Content-Type", "application/json")

	quizUUID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "unable to get uuid"})
		return
	}

	logger = logger.With(slog.String("quizUUID", quizUUID.String()))
	ctx = logging.WithLogger(ctx, logger)

	runnerKey, err := strconv.ParseUint(chi.URLParam(r, "runnerKey"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.ErrMalformedRequest)
		return
	}
	logger = logger.With(slog.Uint64("runner_key", runnerKey))

	tr, ok := c.manager.Get(runnerKey)
	if !ok {
		w.WriteHeader(http.StatusLocked)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: testrunner.ErrRunnerInactive.Error()})
		return
	}

	quizC, err := tr.Get(quizUUID)
	if err != nil {
		logger.Error("couldn't retrive quiz", slog.String("Error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: err.Error()})
		return
	}

	etag := strconv.FormatUint(quizC.Checksum, 10)
	deadline := tr.Deadline()

	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", fmt.Sprintf("private, max-age=%d, must-revalidate", int(math.Round(time.Until(deadline).Hours()))))

	if match := r.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified) // caching goes brr
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contracts.GetQuizFromRunningResponse{
		Quiz: quiz.Quiz{
			Meta:    quizC.Meta,
			Title:   quizC.Title,
			Body:    quizC.Body,
			UUID:    quizC.UUID,
			Options: quizC.Options,
		},
	})
}

func (c *chiService) ListTests(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Failed to parse form data"})
		return
	}

	page, size, jerr := httputils.Page(r)
	if jerr != nil {
		logger.Error("couldn't retrieve page/size from request",
			slog.String("Error", jerr.Error()),
		)
		jerr.Encode(w)
		return
	}

	logger.With(slog.Int("page", page), slog.Int("size", size))

	claims, ok := (r.Context().Value("claims")).(jwtutils.Claims)
	if !ok {
		logger.Error("Failed to retrieve jwt claims")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Unable to retrieve jwt's claims"})
		return
	}

	var tests []models.Test
	var total int
	if claims.IsTeacher {
		logger.Debug("got teacher request")
		tests, total, err = c.testRepo.ListTestsAdvanced(ctx, page, size)
		logger.Debug("got advanced tests info", slog.Any("tests", tests))
	} else {
		logger.Debug("got pupil request", slog.Any("claims", claims))
		tests, total, err = c.testRepo.ListTests(ctx, page, size)
	}

	if err != nil {
		logger.Error("couldn't select tests for listing", slog.String("Error", err.Error()))
		if gone, ok := errors.AsType[carefulness.Gone](err); ok {
			gone.Encode(w)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contracts.ListTestsResponse{
		Tests: tests,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// ExportTest implements [Service].
func (c *chiService) ExportTest(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	best, jerr := httputils.BestAccept(r.Header.Get("Accept"),
		httputils.AvailableArchiveMimes...,
	)
	if jerr != nil || best == "" {
		jerr.Encode(w)
		return
	}
	w.Header().Set("Content-Type", best)

	testUUID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		logger.Error("couldn't parse UUID in url",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't parse UUID in url"})
		return
	}

	logger = logger.With(
		slog.String("test_uuid", testUUID.String()),
	)
	ctx = logging.WithLogger(ctx, logger)

	test, err := c.testRepo.Test(ctx, testUUID)
	if err != nil {
		logger.Error("couldn't select test",
			slog.String("Error", err.Error()),
		)
		if gone, ok := errors.AsType[carefulness.Gone](err); ok {
			gone.Encode(w)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Manifest := exportutils.Manifest{
		Kind: exportutils.Test,
		UUID: testUUID,
		Name: test.Name,
		Quizzes: lo.Map(test.Quizzes, func(quiz models.Quiz, _ int) exportutils.YamlQuiz {
			p, _ := c.cfg.TurnToRel(quiz.Path)
			return exportutils.YamlQuiz{
				UUID:     quiz.UUID,
				Checksum: string(quiz.Checksum.String()),
				Path:     p,
			}
		}),
	}

	out, err := yaml.Marshal(Manifest)
	if err != nil {
		logger.Error("couldn't marshal test yaml",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't create the YAML file for tests"})
		return
	}

	var exporter exportutils.Exporter
	switch best {
	case httputils.Yaml:
		w.Write(out)
		w.WriteHeader(http.StatusOK)
		return
	case httputils.Zip:
		exporter = exportutils.NewZipExporter()
	case httputils.Tar:
		exporter = exportutils.NewTarExporter()
	case httputils.GZip:
		exporter = exportutils.NewGzipExporter()

	default:
		carefulness.ErrUnnacaptable.Encode(w)
		return
	}

	var files []exportutils.ExportFile
	for _, quiz := range test.Quizzes {
		uuid := quiz.UUID
		path, err := c.quizRepo.QuizPath(ctx, uuid)
		if err != nil {
			logger.Error("couldn't get path", slog.String("uuid", uuid.String()), slog.String("error", err.Error()))
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(carefulness.JSONError{Err: fmt.Sprintf("%v not found", uuid)})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(carefulness.JSONError{Err: fmt.Sprintf("unable to process %v", uuid)})
			}
			return
		}
		fi, err := os.Stat(path)
		if err != nil {
			logger.Error("quiz file not accessible", slog.String("path", path), slog.String("Error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		p, err := c.cfg.TurnToRel(path)
		if err != nil {
			logger.Error("unable to turn quiz path to rel", slog.String("Error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(carefulness.JSONError{Err: fmt.Sprintf("unable to turn quiz path to relative")})
			return
		}
		files = append(files, exportutils.ExportFile{UUID: uuid.String(), Path: path, Name: p, FileInfo: fi})
	}

	tmpDir, err := os.MkdirTemp(filepath.Dir(c.cfg.QuizPath), ".import-")
	if err != nil {
		logger.Error("failed to create tmpDir",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't create temp dir for new quiz bank"})
		return
	}
	defer os.RemoveAll(tmpDir)

	f, err := os.CreateTemp(tmpDir, fmt.Sprintf("%v-manifest-*.yaml", test.UUID.String()))
	if err != nil {
		logger.Error("unable to create temp file", slog.String("Error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: fmt.Sprintf("unable to create temp file")})
		return
	}
	defer f.Close()

	_, err = io.Copy(f, bytes.NewReader(out))
	if err != nil {
		logger.Error("unable to copy to temp file", slog.String("Error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: fmt.Sprintf("unable to copy manifest to temp file")})
		return
	}

	info, err := f.Stat()
	if err != nil {
		logger.Error("unable to get info about temp file", slog.String("Error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: fmt.Sprintf("unable to get info about temp file")})
		return
	}

	var buf bytes.Buffer
	err = exporter.Export(ctx, &buf, append(files, exportutils.ExportFile{
		FileInfo: info,
		Path:     filepath.Join(tmpDir, info.Name()),
		UUID:     test.UUID.String(),
		Name:     info.Name(),
	}))
	if err != nil {
		logger.Error("couldn't export test", slog.String("Error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "error while writing to archive"})
		return
	}
	w.Header().Set("Content-Type", best)
	buf.WriteTo(w)
}

// ImportTest implements [Service].
func (c *chiService) ImportTest(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Add("content-type", "application/json")

	if err := r.ParseMultipartForm(8 << 20); err != nil {
		logger.Error("archive is too big", slog.String("Error", err.Error()))
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "archive is too big!"})
		return
	}

	archiveParts, handler, err := r.FormFile("imported")
	if err != nil {
		logger.Error("couldn't get file", slog.String("Error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "unable to parse file"})
		return
	}
	defer archiveParts.Close()

	contentType, _, err := mime.ParseMediaType(handler.Header.Get("Content-Type"))
	if err != nil {
		logger.Error("couldn't parse MIME type", slog.String("Error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "unable to parse MIME"})
		return
	}
	if !(httputils.ValidArchive(contentType)) {
		carefulness.ErrUnnacaptable.Encode(w)
		return
	}

	logger = logger.With(
		slog.String("file_name", handler.Filename),
		slog.String("mime", contentType),
		slog.Int64("size", handler.Size),
	)
	ctx = logging.WithLogger(ctx, logger)

	tmpDir, err := os.MkdirTemp(filepath.Dir(c.cfg.QuizPath), ".import-")
	if err != nil {
		logger.Error("failed to create tmpDir",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't create temp dir for new quiz bank"})
		return
	}
	defer os.RemoveAll(tmpDir)

	var importer importutils.Importer
	switch contentType {
	case httputils.Zip:
		importer = importutils.NewZipImporter()
	case httputils.Tar:
		importer = importutils.NewTarImporter()
	case httputils.GZip:
		importer = importutils.NewGzipImporter()
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	staged, test, jerr := importer.Import(ctx, archiveParts, handler.Size, tmpDir, c.cfg.TurnToAbs)
	if jerr != nil {
		logger.Error("couldn't unarchive an archive", slog.String("Error", jerr.Error()))
		jerr.Encode(w)
		return
	}

	err = c.commitTestImport(ctx, staged, test)
	if err != nil {
		logger.Error("couldn't commit the import", slog.String("Error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't commit the import"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *chiService) GetRunningQuizzesUUIDs(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	w.Header().Set("Content-Type", "application/json")

	runnerKey, err := strconv.ParseUint(chi.URLParam(r, "runnerKey"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.ErrMalformedRequest)
		return
	}
	logger = logger.With(slog.Uint64("runner_key", runnerKey))

	tr, ok := c.manager.Get(runnerKey)
	if !ok {
		w.WriteHeader(http.StatusLocked)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: testrunner.ErrRunnerInactive.Error()})
		return
	}

	checksum := tr.Checksum()
	if checksum == 0 {
		w.WriteHeader(http.StatusLocked)
		return
	}

	etag := strconv.FormatUint(checksum, 10)
	deadline := tr.Deadline()

	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", fmt.Sprintf("private, max-age=%d, must-revalidate", int(math.Round(time.Until(deadline).Hours()))))

	if match := r.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified) // caching goes brrrrr
		return
	}

	uuids := tr.GetAll()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contracts.GetQuizzesUUIDs{
		UUIDs: uuids,
	})
}

func (c *chiService) RunningInfo(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	w.Header().Set("Content-Type", "application/json")
	ctx := logging.WithLogger(r.Context(), logger)

	runners := c.manager.AllRunners()
	testUUIDs := make(uuid.UUIDs, 0, len(runners))
	for _, testUUID := range runners {
		testUUIDs = append(testUUIDs, testUUID)
	}

	tests, err := c.testRepo.Tests(ctx, testUUIDs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't select running test info"})
		return
	}

	testsByUUID := make(map[uuid.UUID]models.Test, len(tests))
	for _, test := range tests {
		testsByUUID[test.UUID] = test
	}

	testInfos := make([]contracts.RunningInfo, 0, len(runners))
	for key, testUUID := range runners {
		tr, ok := c.manager.Get(key)
		if !ok {
			continue
		}
		test, ok := testsByUUID[testUUID]
		if !ok {
			continue // me when I delete test from db mid-run
		}

		testInfos = append(testInfos, contracts.RunningInfo{
			Key:      key,
			Deadline: tr.Deadline(),
			IsPaused: tr.IsPaused(),
			Test:     test,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contracts.RunningInfoResponse{
		TestInfos: testInfos,
	})
}

func (c *chiService) GetDeadline(key uint64) (time.Time, error) {
	tr, ok := c.manager.Get(key)
	if !ok {
		return time.Time{}, testrunner.ErrRunnerInactive
	}

	return tr.Deadline(), nil
}

func (c *chiService) GetTestUUID(key uint64) uuid.UUID {
	tr, ok := c.manager.Get(key)
	if !ok {
		return uuid.Nil
	}

	return tr.Test()
}
