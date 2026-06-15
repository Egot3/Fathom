package testutils

import (
	"context"
	"fmt"
	"log"

	"github.com/egot3/fathom/internal/testutils/migrations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

func RunMigrations(ctx context.Context, db *bun.DB) error {
	migrator := migrate.NewMigrator(db, migrations.Migrations)

	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("Migration init failed")
	}

	for {
		log.Printf("All migrations: %d", len(migrations.Migrations.Sorted()))
		log.Printf("Unaplied migrations: %d", len(migrations.Migrations.Sorted().Unapplied()))
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
