package logging

import (
	"context"
	"log/slog"
)

type key struct{}

var fallback = NewLogger([]string{"charmLog"}, "info")

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, key{}, logger)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(key{}).(*slog.Logger); ok {
		return logger
	}

	return fallback
}
