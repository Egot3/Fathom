package migrations

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/egot3/fathom/internal/logging"
	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/migrate"
)

//go:embed sqlite/*.sql
var sqliteMigrations embed.FS

//go:embed postgres/*.sql
var postgresMigrations embed.FS

func New(ctx context.Context, dialectName dialect.Name) (*migrate.Migrations, error) {
	logger := logging.LoggerFromContext(ctx).With(
		slog.String("layer", "migration"),
	)
	ctx = logging.WithLogger(ctx, logger)

	migrations := migrate.NewMigrations()

	var fsys embed.FS
	switch dialectName {
	case dialect.PG:
		fsys = postgresMigrations
	case dialect.SQLite:
		fsys = sqliteMigrations
	default:
		return nil, fmt.Errorf("invalid dialect: %v", dialectName)
	}

	if err := migrations.Discover(fsys); err != nil {
		return nil, err
	}
	return migrations, nil
}
