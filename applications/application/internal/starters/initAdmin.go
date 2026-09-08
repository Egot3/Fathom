package starters

import (
	"context"
	"log/slog"
	"os"

	"github.com/egot3/fathom/internal/config"
	"github.com/egot3/fathom/internal/models"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

func initAdmin(i do.Injector) {
	cfg := do.MustInvoke[*config.Config](i)
	logger := do.MustInvoke[*slog.Logger](i)
	db := do.MustInvoke[*bun.DB](i)

	if cfg.InitAdminPassword != "" && cfg.InitAdminUsername != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(cfg.InitAdminPassword), bcrypt.DefaultCost)
		if err != nil {
			logger.Error("Couldn't create init teacher", slog.String("Error", err.Error()))
			os.Exit(1)
		}
		_, err = db.NewInsert().On("CONFLICT DO UPDATE").Model(&models.User{
			Nickname:     cfg.InitAdminUsername,
			PasswordHash: passwordHash,
			IsTeacher:    true,
		}).Exec(context.Background())
		if err != nil {
			logger.Error("Couldn't create init teacher: %v", slog.String("Error", err.Error()))
			os.Exit(1)
		}
	}
}
