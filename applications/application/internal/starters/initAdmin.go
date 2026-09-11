package starters

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/egot3/fathom/internal/config"
	"github.com/egot3/fathom/internal/models"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

func initAdminErrorable(i do.Injector) error {
	cfg := do.MustInvoke[*config.Config](i)
	db := do.MustInvoke[*bun.DB](i)

	if cfg.InitAdminPassword == "" || cfg.InitAdminUsername == "" {
		return nil
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cfg.InitAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.Join(errors.New("couldn't hash password"), err)
	}
	_, err = db.NewInsert().On("CONFLICT DO UPDATE").Model(&models.User{
		Nickname:     cfg.InitAdminUsername,
		PasswordHash: passwordHash,
		IsTeacher:    true,
	}).Exec(context.Background())
	if err != nil {
		return errors.Join(errors.New("couldn't insert init teacher"), err)
	}

	return nil
}

func initAdmin(i do.Injector) {
	logger := do.MustInvoke[*slog.Logger](i)

	if err := initAdminErrorable(i); err != nil {
		logger.Error("Unable to init teachers", slog.String("Error", err.Error()))
		os.Exit(1)
	}
}
