package migrations

import (
	"embed"
	"log"

	"github.com/uptrace/bun/migrate"
)

//go:embed *.sql
var sqlMigrations embed.FS

var Migrations = migrate.NewMigrations()

func init() {
	entries, err := sqlMigrations.ReadDir(".")
	if err != nil {
		panic(err)
	}
	log.Println("Embedded files:")
	for _, e := range entries {
		log.Printf("- %s", e.Name())
	}

	if err := Migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}
	log.Printf("Discovered %d migration groups", len(Migrations.Sorted()))
}
