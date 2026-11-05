package testutils

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/egot3/fathom/internal/config"
	"github.com/samber/do/v2"
)

func WithConfig(tb testing.TB, i do.Injector) {
	cfg := &config.Config{QuizPath: filepath.Join(tb.TempDir(), "quizzes")}
	do.ProvideValue(i, cfg)
}

func WithLogger(tb testing.TB, i do.Injector) {
	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)

	do.ProvideValue(i, slog.New(logger))
}
