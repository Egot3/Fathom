package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/egot3/fathom/internal/config"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func InitDB(i do.Injector) (*bun.DB, error) {
	cfg := do.MustInvoke[config.Config](i)

	var dsn string
	var sqldb *sql.DB = nil
	var err error
	var DB *bun.DB = nil
	if cfg.Database.Sqlite.Used {
		dsn = cfg.Database.Sqlite.Path

		sqldb, err = sql.Open(sqliteshim.ShimName, fmt.Sprintf("file:%v", dsn))
		if err != nil {
			return nil, err
		}
		DB = bun.NewDB(sqldb, sqlitedialect.New())
	}
	if cfg.Database.Postgres.Used {
		dsn = fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
			cfg.Database.Postgres.User,
			cfg.Database.Postgres.Password,
			cfg.Database.Postgres.Host,
			cfg.Database.Postgres.Port,
			cfg.Database.Postgres.DbName)

		sqldb = sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
		DB = bun.NewDB(sqldb, pgdialect.New())
	}

	if sqldb == nil {
		return nil, fmt.Errorf("db was not defined, please check config")
	}

	for i := range 5 {
		if err := DB.Ping(); err != nil {
			log.Printf("Try %d: Pings didn't pong: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}
		break
	} //почему бы и нет

	if err := DB.Ping(); err != nil {
		log.Printf("\nNo db?\n")
		return nil, err
	}

	sqldb.SetMaxOpenConns(50)
	sqldb.SetMaxIdleConns(20)

	log.Printf("DB UP")
	return DB, nil
}
