package starters

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/egot3/fathom/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
)

type Starter interface {
	Serve() error
	Shutdown(ctx context.Context) error
}

type httpStarter struct {
	logger *slog.Logger
	server *http.Server
}

func (h httpStarter) Serve() error {
	h.logger.Info("started serving!", slog.String("address", h.server.Addr))
	err := h.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (h httpStarter) Shutdown(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}

func newHTTPStarter(i do.Injector) (Starter, error) {
	cfg := do.MustInvoke[*config.Config](i)
	handler := do.MustInvoke[chi.Router](i)
	logger := do.MustInvoke[*slog.Logger](i)

	return httpStarter{
		logger: logger,
		server: &http.Server{
			Addr:    ":" + cfg.ServerPort,
			Handler: handler,
		},
	}, nil
}
