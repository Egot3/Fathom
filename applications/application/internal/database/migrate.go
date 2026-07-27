package database

import (
	"context"
	"fmt"
	"log"

	"github.com/egot3/fathom/internal/database/migrations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

func RunMigrations(ctx context.Context, db *bun.DB) error {
	migrations, err := migrations.New(ctx, db.Dialect().Name())
	if err != nil {
		return fmt.Errorf("build migrations: %w", err)
	}

	migrator := migrate.NewMigrator(db, migrations)
	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("Migration init failed")
	}

	for {
		log.Printf("All migrations: %d", len(migrations.Sorted()))
		log.Printf("Unaplied migrations: %d", len(migrations.Sorted().Unapplied()))
		group, err := migrator.Migrate(ctx)
		if err != nil {
			return fmt.Errorf("migration failed: %v", err)
		}

		if group.IsZero() {
			log.Println("all migrations applied")
			break
		}
		log.Printf("Migrated to %s", group)

	}
	return nil

}
