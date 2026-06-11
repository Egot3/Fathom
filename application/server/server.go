package server

import (
	"log/slog"
	"net/http"

	"github.com/egot3/fathom/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
)

func ChiServer(i do.Injector) (chi.Router, error) {
	r := chi.NewRouter()

	r.Use(middleware.AttachLogger(do.MustInvoke[*slog.Logger](i)))
	r.Use(middleware.TraceAttacher)

	r.Method("GET", "/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthy"))
		w.WriteHeader(http.StatusOK)
	}))

	return r, nil
}
