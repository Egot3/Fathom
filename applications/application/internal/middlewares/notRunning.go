package middlewares

import (
	"log/slog"
	"net/http"
	"slices"

	"github.com/egot3/fathom/internal/logging"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// requires "test_uuid" in URL
func TestNotRunning(uuidGetter func() uuid.UUID) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logging.LoggerFromContext(r.Context())
			logger = logger.With(slog.String("layer", "middleware"))

			if chi.URLParam(r, "test_uuid") == uuidGetter().String() {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Running(uuidGetter func() uuid.UUID) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logging.LoggerFromContext(r.Context())
			logger = logger.With(slog.String("layer", "middleware"))

			logger.Debug("checking if test uuid is running")
			if uuidGetter() != uuid.Nil {
				next.ServeHTTP(w, r)
				return
			}

			logger.Debug("test uuid is not running")
			w.WriteHeader(http.StatusLocked)
		})
	}
}

// requires quiz_uuid as URL param
func QuizNotRunning(uuidsGetter func() uuid.UUIDs) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logging.LoggerFromContext(r.Context())
			logger = logger.With(slog.String("layer", "middleware"))

			logger.Debug("checking if quiz uuid is in running")
			if slices.Contains(uuidsGetter().Strings(), chi.URLParam(r, "quiz_uuid")) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			logger.Debug("it's'n't running, proceed")
			next.ServeHTTP(w, r)
		})
	}
}
