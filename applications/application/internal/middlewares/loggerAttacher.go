package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/egot3/fathom/internal/logging"
)

func AttachLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rctx := r.Context()

			ctx := logging.WithLogger(rctx, logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
