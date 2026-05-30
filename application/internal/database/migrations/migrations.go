package migrations

import (
	"context"
	"embed"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/migrate"
)

//go:embed sqlite/*.sql
var sqliteMigrations embed.FS

//go:ember postgres/*.sql
var postgresMigrations embed.FS

var Migrations = migrate.NewMigrations()

func init() {

	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		switch db.Dialect().Name() {
		case dialect.PG:
			entries, err := postgresMigrations.ReadDir("./postgres")
			if err != nil {
				panic(err)
			}
			log.Println("Embedded files:")
			for _, e := range entries {
				log.Printf("- %s", e.Name())
			}

			if err := Migrations.Discover(postgresMigrations); err != nil {
				panic(err)
			}
		case dialect.SQLite:
			entries, err := sqliteMigrations.ReadDir("./sqlite")
			if err != nil {
				panic(err)
			}
			log.Println("Embedded files:")
			for _, e := range entries {
				log.Printf("- %s", e.Name())
			}

			if err := Migrations.Discover(sqliteMigrations); err != nil {
				panic(err)
			}
		default:
			panic("Invalid db type")
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		return nil
	})

	log.Printf("Discovered %d migration groups", len(Migrations.Sorted()))
}
