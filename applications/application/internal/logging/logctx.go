package logging

import (
	"context"
	"log/slog"
	"os"
)

type key struct{}

var fallback = func() *slog.Logger {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return slog.New(jsonHandler)
}()

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, key{}, logger)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(key{}).(*slog.Logger); ok {
		return logger
	}

	return fallback
}
