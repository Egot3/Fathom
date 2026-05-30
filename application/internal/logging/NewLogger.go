package logging

import (
	"log/slog"
	"os"
	"slices"

	charmlog "github.com/charmbracelet/log"
	"github.com/charmbracelet/x/term"
	slogmulti "github.com/samber/slog-multi"
)

func NewLogger(loggers []string, levelStr string) *slog.Logger {
	handlers := []slog.Handler{}

	var level slog.Level
	switch levelStr {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	if slices.Contains(loggers, "slog") {
		jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
		handlers = append(handlers, jsonHandler)
	}
	if slices.Contains(loggers, "charmLog") {
		if term.IsTerminal(os.Stderr.Fd()) {
			charmHandler := charmlog.New(os.Stderr)
			charmHandler.SetLevel(charmlog.Level(level))
			handlers = append(handlers, charmHandler)
		}
	}

	return slog.New(slogmulti.Fanout(handlers...))
}
