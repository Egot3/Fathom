package middlewares

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/egot3/fathom/internal/logging"
	"github.com/google/uuid"
)

type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (rwi *responseWriterInterceptor) WriteHeader(statusCode int) {
	rwi.statusCode = statusCode
	rwi.ResponseWriter.WriteHeader(statusCode)
}

func (rwi *responseWriterInterceptor) Write(b []byte) (int, error) {
	if rwi.statusCode == 0 {
		rwi.statusCode = http.StatusOK
	}
	return rwi.ResponseWriter.Write(b)
}

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

		interceptor := &responseWriterInterceptor{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(interceptor, r.WithContext(ctx))

		code := interceptor.statusCode
		elapsed := time.Since(start)
		logger.InfoContext(ctx, "done",
			slog.Duration("elapsed", elapsed),
			slog.Int("code", code),
		)
	})
}
