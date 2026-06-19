package migrations

import (
	"embed"
	"fmt"

	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/migrate"
)

//go:embed sqlite/*.sql
var sqliteMigrations embed.FS

//go:embed postgres/*.sql
var postgresMigrations embed.FS

func New(dialectName dialect.Name) (*migrate.Migrations, error) {
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
