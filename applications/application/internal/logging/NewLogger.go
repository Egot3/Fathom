package logging

import (
	"fmt"
	"log/slog"
	"os"
	"slices"

	charmlog "github.com/charmbracelet/log"
	"github.com/charmbracelet/x/term"
	"github.com/egot3/fathom/internal/config"
	"github.com/samber/do/v2"
	slogmulti "github.com/samber/slog-multi"
)

// loggers []string, levelStr string
func NewLogger(i do.Injector) (*slog.Logger, error) {
	handlers := []slog.Handler{}
	cfg := do.MustInvoke[config.Config](i)

	var level slog.Level
	switch cfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		return nil, fmt.Errorf("no Log level")
	}

	if len(cfg.LogSinks) > 2 {
		return nil, fmt.Errorf("Can't have other loggers than slog and charm")
	}
	if slices.Contains(cfg.LogSinks, "slog") {
		jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
		handlers = append(handlers, jsonHandler)
	}
	if slices.Contains(cfg.LogSinks, "charmLog") {
		if term.IsTerminal(os.Stderr.Fd()) {
			charmHandler := charmlog.New(os.Stderr)
			charmHandler.SetLevel(charmlog.Level(level))
			handlers = append(handlers, charmHandler)
		}
	}

	return slog.New(slogmulti.Fanout(handlers...)), nil
}
