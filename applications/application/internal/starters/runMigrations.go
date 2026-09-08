package starters

import (
	"context"
	"log/slog"

	"github.com/egot3/fathom/internal/database"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

func runMigrations(i do.Injector) {
	db := do.MustInvoke[*bun.DB](i)
	logger := do.MustInvoke[*slog.Logger](i)

	if err := database.RunMigrations(context.Background(), db); err != nil {
		logger.Error("Fatal migration error", slog.String("Error", err.Error()))
		panic(err)
	}
}
