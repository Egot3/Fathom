package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/egot3/fathom/internal/starters"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

func main() {
	i := starters.GenerateInjector()
	logger := do.MustInvoke[*slog.Logger](i)
	starter := do.MustInvoke[starters.Starter](i)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serveErrCh := make(chan error, 1)
	go func() {
		serveErrCh <- starter.Serve()
	}()

	select {
	case err := <-serveErrCh:
		if err != nil {
			logger.Error("server failed to start", slog.String("error", err.Error()))
			os.Exit(1)
		}
	case <-ctx.Done():
		logger.Info("shutdown signal received, ceasing-and-draining")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := starter.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", slog.String("error", err.Error()))
	}

	db := do.MustInvoke[*bun.DB](i)
	if err := db.Close(); err != nil {
		logger.Error("db close failed", slog.String("error", err.Error()))
	}
}
