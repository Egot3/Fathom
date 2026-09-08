package main

import (
	"log/slog"
	"os"

	"github.com/egot3/fathom/internal/starters"
	"github.com/samber/do/v2"
)

func main() {
	i := starters.GenerateInjector()

	logger := do.MustInvoke[*slog.Logger](i)
	starter := do.MustInvoke[starters.Starter](i)

	if err := starter.Serve(); err != nil {
		logger.Info("Server execution finished", slog.String("Error", err.Error()))
		os.Exit(0)
	}
}
