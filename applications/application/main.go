package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	jwtutils "github.com/egot3/fathom/internal/JWTutils"
	"github.com/egot3/fathom/internal/starters"
	testrunner "github.com/egot3/fathom/internal/testRunner"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

func main() {
	i := starters.GenerateInjector()
	logger := do.MustInvoke[*slog.Logger](i)
	starter := do.MustInvoke[starters.Starter](i)
	err := jwtutils.Load()
	if err != nil {
		logger.Error("bad JWT_SECRET", slog.String("Error", err.Error()))
		os.Exit(1)
	}
	serveErrCh := make(chan error, 1)
	go func() {
		serveErrCh <- starter.Serve()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serveErrCh:
		if err != nil {
			logger.Error("server failed to start", slog.String("error", err.Error()))
			os.Exit(1)
		}
	case sig := <-sigCh:
		logger.Info("shutdown signal received, ceasing-and-draining", slog.String("signal", sig.String()))
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := do.ShutdownWithContext[starters.Starter](shutdownCtx, i); err != nil {
		logger.Error("graceful shutdown failed on starter", slog.String("error", err.Error()))
	} else {
		logger.Info("starter shot down")
	}

	if err := do.ShutdownWithContext[*testrunner.Manager](shutdownCtx, i); err != nil {
		logger.Error("graceful shutdown failed on manager", slog.String("error", err.Error()))
	} else {
		logger.Info("manager stopped")
	}

	db := do.MustInvoke[*bun.DB](i)
	if err := db.Close(); err != nil {
		logger.Error("db close failed", slog.String("error", err.Error()))
	} else {
		logger.Info("db connection closed")
	}

	logger.Info("shutdown finished", slog.String("goodbye", "have a good one"))
}
