package hooks

import (
	"context"
	"log/slog"
	"time"

	"github.com/egot3/fathom/internal/logging"
	"github.com/uptrace/bun"
)

type SlogQueryHook struct{}

func (h *SlogQueryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	logging.LoggerFromContext(ctx).InfoContext(ctx, "starting db request",
		slog.String("operation", event.Operation()),
		slog.String("query", event.Query),
	)
	return ctx
}

func (h *SlogQueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	elapsed := time.Since(event.StartTime)
	logger := logging.LoggerFromContext(ctx).With(
		slog.String("operation", event.Operation()),
		slog.Duration("elapse", elapsed),
	)

	if event.Err != nil {
		logger.ErrorContext(ctx, "db request ffailed", slog.String("err", event.Err.Error()))
		return
	}

	logger.InfoContext(ctx, "db request succeeded")
}
