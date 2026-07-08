package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/egot3/fathom/internal/logging"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// requires "test_uuid" in URL
func NotRunning(uuidGetter func() uuid.UUID) func(http.Handler) http.Handler {
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
