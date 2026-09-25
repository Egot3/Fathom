package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/contracts"
	"github.com/egot3/fathom/internal/httputils"
	"github.com/egot3/fathom/internal/logging"
	"github.com/egot3/fathom/internal/quiz"
	testrunner "github.com/egot3/fathom/internal/testRunner"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// GetAnswer implements [Service].
func (c *chiService) GetAnswer(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)
	w.Header().Set("Content-Type", "application/json")

	groupUUID, err := uuid.Parse(chi.URLParam(r, "group_uuid"))
	if err != nil {
		logger.Error("couldn't parse testUUID in url",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userUUID, err := uuid.Parse(chi.URLParam(r, "user_uuid"))
	if err != nil {
		logger.Error("couldn't parse userUUID in url",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	quizUUID, err := uuid.Parse(chi.URLParam(r, "quiz_uuid"))
	if err != nil {
		logger.Error("couldn't parse quizUUID in url",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	testUUID, err := uuid.Parse(chi.URLParam(r, "test_uuid"))
	if err != nil {
		logger.Error("couldn't parse testUUID in url",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger = logger.With(
		slog.String("group_uuid", groupUUID.String()),
		slog.String("user_uuid", userUUID.String()),
		slog.String("test_uuid", testUUID.String()),
		slog.String("quiz_uuid", quizUUID.String()),
	)
	ctx = logging.WithLogger(ctx, logger)

	answer, err := c.answerRepo.Answer(ctx, userUUID, testUUID, groupUUID, quizUUID)
	if err != nil {
		logger.Error("couldn't get an answer",
			slog.String("Error", err.Error()),
		)

		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Requested answer is not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't get an answer because of unknown error"})
		return
	}
	answerJSON := quiz.QuizAnswers{}
	err = json.Unmarshal([]byte(answer), &answerJSON)
	if err != nil {
		logger.Error("couldn't parse the user's answer",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "unable to parse user's answer"})
		return
	}

	correct, err := c.quizRepo.CorrectAnswer(ctx, quizUUID)
	if err != nil {
		logger.Error("couldn't get the right answer",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "unable to get correct answer"})
		return
	}
	correctJSON := quiz.QuizAnswers{}
	err = json.Unmarshal([]byte(correct), &correctJSON)
	if err != nil {
		logger.Error("couldn't parse the right answer",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "unable to parse correct answer"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contracts.AnswerResponse{
		Answer: contracts.Answer{
			GroupUUID: groupUUID,
			TestUUID:  testUUID,
			UserUUID:  userUUID,
			QuizUUID:  quizUUID,
			Chosen:    answerJSON,
			Correct:   correctJSON,
		},
	})
}

// PostAnswer implements [Service].
func (c *chiService) PostAnswer(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)
	w.Header().Set("Content-Type", "application/json")

	quizUUID, err := uuid.Parse(chi.URLParam(r, "quiz_uuid"))
	if err != nil {
		logger.Error("couldn't parse testUUID in url",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	groupUUID, err := uuid.Parse(chi.URLParam(r, "group_uuid"))
	if err != nil {
		logger.Error("couldn't parse testUUID in url",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userUUID, err := uuid.Parse(chi.URLParam(r, "user_uuid"))
	if err != nil {
		logger.Error("couldn't parse groupUUID in url",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger = logger.With(
		slog.String("quiz_uuid", quizUUID.String()),
	)
	ctx = logging.WithLogger(ctx, logger)

	var req contracts.PostAnswerRequest
	jerr := httputils.ParseJSON(r.Body, &req)
	if jerr != nil {
		logger.Error("Failed to parse body",
			slog.String("Error", jerr.Error()),
		)
		jerr.Encode(w)
		return
	}

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

	q, err := tr.Get(quizUUID)
	if err != nil {
		logger.Error("couldn't get quiz",
			slog.String("Error", err.Error()),
		)
		switch {
		case errors.Is(err, testrunner.ErrQuizNotCached):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, testrunner.ErrRunnerInactive):
			w.WriteHeader(http.StatusLocked)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(carefulness.JSONError{Err: err.Error()})
		return
	}

	score := q.EvaluateScore(req.Value)

	answered, err := json.Marshal(req.Value)
	if err != nil {
		logger.Error("couldn't marshal answer to json(how?)",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusTeapot)
		return
	}
	err = c.answerRepo.SetAnswer(ctx, tr.Test(), groupUUID, userUUID, quizUUID, string(answered), score)
	if err != nil {
		logger.Error("couldn't set an answer",
			slog.String("Error", err.Error()),
		)

		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't set an answer because of unknown error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Totalize implements [Service].
func (c *chiService) Totalize(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)
	w.Header().Set("Content-Type", "application/json")

	groupUUID, err := uuid.Parse(chi.URLParam(r, "group_uuid"))
	if err != nil {
		logger.Error("couldn't parse groupUUID in url",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userUUID, err := uuid.Parse(chi.URLParam(r, "user_uuid"))
	if err != nil {
		logger.Error("couldn't parse userUUID in url",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	currentTestUUID := tr.Test()
	logger = logger.With(
		slog.String("group_uuid", groupUUID.String()),
		slog.String("user_uuid", userUUID.String()),
		slog.String("currently_running_test_uuid", currentTestUUID.String()),
	)
	ctx = logging.WithLogger(ctx, logger)

	err = c.answerRepo.Totalize(ctx, userUUID, currentTestUUID, groupUUID)
	if err != nil {
		logger.Error("couldn't totalize",
			slog.String("Error", err.Error()),
		)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't totalize because of unknown error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *chiService) ListTotals(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)
	w.Header().Set("Content-Type", "application/json")

	page, size, jerr := httputils.Page(r)
	if jerr != nil {
		logger.Error("couldn't retrieve page/size from request",
			slog.String("Error", jerr.Error()),
		)
		jerr.Encode(w)
		return
	}

	logger = logger.With(slog.Int("page", page), slog.Int("size", size))

	var err error

	groupUUID := uuid.Nil
	if raw := chi.URLParam(r, "group_uuid"); raw != "all" {
		groupUUID, err = uuid.Parse(raw)
		if err != nil {
			logger.Error("couldn't parse groupUUID in url",
				slog.String("Error", err.Error()),
			)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	userUUID := uuid.Nil
	if raw := chi.URLParam(r, "user_uuid"); raw != "all" {
		userUUID, err = uuid.Parse(raw)
		if err != nil {
			logger.Error("couldn't parse userUUID in url",
				slog.String("Error", err.Error()),
			)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	testUUID := uuid.Nil
	if raw := chi.URLParam(r, "test_uuid"); raw != "all" {
		testUUID, err = uuid.Parse(raw)
		if err != nil {
			logger.Error("couldn't parse testUUID in url",
				slog.String("Error", err.Error()),
			)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	logger = logger.With(
		slog.String("test_uuid", testUUID.String()),
		slog.String("user_uuid", userUUID.String()),
		slog.String("group_uuid", groupUUID.String()),
	)

	totals, total, err := c.answerRepo.ListTotals(ctx, page, size,
		userUUID, testUUID, groupUUID)
	if err != nil {
		logger.Error("couldn't get an answer",
			slog.String("Error", err.Error()),
		)

		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Err: "there is no totals to list"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't get an answer because of unknown error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contracts.TotalsResponse{
		Totals: totals,
		Total:  total,
		Page:   page,
		Size:   size,
	})
}

func (c *chiService) ListUserAnswer(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)
	w.Header().Set("Content-Type", "application/json")

	groupUUID, err := uuid.Parse(chi.URLParam(r, "group_uuid"))
	if err != nil {
		logger.Error("couldn't parse groupUUID in url",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userUUID, err := uuid.Parse(chi.URLParam(r, "user_uuid"))
	if err != nil {
		logger.Error("couldn't parse groupUUID in url",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	testUUID, err := uuid.Parse(chi.URLParam(r, "test_uuid"))
	if err != nil {
		logger.Error("couldn't parse testUUID in url",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger = logger.With(
		slog.String("group_uuid", groupUUID.String()),
		slog.String("user_uuid", userUUID.String()),
		slog.String("test_uuid", testUUID.String()),
	)
	ctx = logging.WithLogger(ctx, logger)

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Failed to parse form data"})
		return
	}

	page, size, err := httputils.Page(r)
	if err != nil {
		logger.Error("couldn't retrieve page/size from request",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger = logger.With(slog.Int("page", page), slog.Int("size", size))

	userAnswers, total, err := c.answerRepo.AnswersInTest(ctx, userUUID, testUUID, groupUUID, page, size)
	if err != nil {
		logger.Error("couldn't get an answer",
			slog.String("Error", err.Error()),
		)

		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Err: "there is no totals to list"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "couldn't get an answer because of unknown error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contracts.AnswersResponse{
		Answers: userAnswers,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}
