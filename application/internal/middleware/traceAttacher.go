package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/egot3/fathom/internal/logging"
	"github.com/google/uuid"
)

func TraceAttacher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rctx := r.Context()
		baseLogger := logging.LoggerFromContext(rctx)

		start := time.Now()
		requestUUID := uuid.NewString()

		logger := baseLogger.With(
			slog.String("url", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("traceUUID", requestUUID),
		)
		ctx := logging.WithLogger(rctx, logger)

		next.ServeHTTP(w, r.WithContext(ctx))

		code := r.Response.StatusCode
		elapsed := time.Since(start)
		logger.InfoContext(ctx, "done",
			slog.Duration("elapsed", elapsed),
			slog.Int("code", code),
		)
	})
}
