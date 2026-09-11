package database

import (
	"context"
	"log/slog"
	"os"

	"github.com/egot3/fathom/internal/database/migrations"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

func RunMigrations(i do.Injector) {
	logger := do.MustInvoke[*slog.Logger](i)
	db := do.MustInvoke[*bun.DB](i)

	ctx := context.Background()

	migrations, err := migrations.New(ctx, db.Dialect().Name())
	if err != nil {
		logger.Error("couldn't build migrations", slog.String("Error", err.Error()))
		os.Exit(1)
	}

	migrator := migrate.NewMigrator(db, migrations)
	if err := migrator.Init(ctx); err != nil {
		logger.Error("couldn't init migrations", slog.String("Error", err.Error()))
		os.Exit(1)
	}

	for {
		logger.Debug("migrations detected",
			slog.Int("count", len(migrations.Sorted())),
			slog.Int("unnaplied", len(migrations.Sorted().Unapplied())),
		)
		group, err := migrator.Migrate(ctx)
		if err != nil {
			logger.Error("migration failed", slog.String("Error", err.Error()))
			os.Exit(1)
		}

		if group.IsZero() {
			logger.Debug("all migrations applied")
			break
		}
		logger.Debug("migrated", slog.String("group", group.String()))

	}
}
