package starters

import (
	"net/http"

	"github.com/egot3/fathom/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
)

type Starter interface {
	Serve() error
}

type httpStarter struct {
	port    string // really hate like it looks. uint64 is impossible to conv
	handler http.Handler
}

func (h httpStarter) Serve() error {
	return http.ListenAndServe(":"+h.port, h.handler)
}

func newHTTPStarter(i do.Injector) (Starter, error) {
	cfg := do.MustInvoke[*config.Config](i)
	handler := do.MustInvoke[chi.Router](i)

	return httpStarter{
		port:    cfg.ServerPort,
		handler: handler,
	}, nil
}
