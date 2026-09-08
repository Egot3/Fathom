package database

import (
	"database/sql"
	"fmt"
	"log/slog"
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
	cfg := do.MustInvoke[*config.Config](i)
	logger := do.MustInvoke[*slog.Logger](i)

	var sqldb *sql.DB = nil
	var err error
	var DB *bun.DB = nil

	switch cfg.DatabaseDriver {
	case "postgres":
		sqldb = sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.PostgresURL)))
		DB = bun.NewDB(sqldb, pgdialect.New())
	default:
		sqldb, err = sql.Open(sqliteshim.ShimName, fmt.Sprintf("file:%v", cfg.SqlitePath))
		if err != nil {
			return nil, err
		}
		DB = bun.NewDB(sqldb, sqlitedialect.New())
	}

	if sqldb == nil {
		return nil, fmt.Errorf("db was not defined, please check config")
	}

	for i := range 5 {
		if err := DB.Ping(); err != nil {
			logger.Info("Ping did't pong", slog.Int("count", i+1), slog.String("Error", err.Error()))
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	if err := DB.Ping(); err != nil {
		logger.Info("DB's health couldn't be checked. Our condolences")
		return nil, err
	}

	sqldb.SetMaxOpenConns(50)
	sqldb.SetMaxIdleConns(20)

	return DB, nil
}
