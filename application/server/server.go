package server

import (
	"log/slog"
	"net/http"

	"github.com/egot3/fathom/internal/handler"
	"github.com/egot3/fathom/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
)

func ChiServer(i do.Injector) (chi.Router, error) {
	r := chi.NewRouter()
	svc := do.MustInvoke[handler.Service](i)

	r.Use(middleware.BodySizer)

	r.Method("GET", "/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("healthy"))
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.AttachLogger(do.MustInvoke[*slog.Logger](i)))
		r.Use(middleware.TraceAttacher)

		r.Route("/user", func(r chi.Router) {

			r.Group(func(r chi.Router) {
				r.With(middleware.ParseUUID).Get("/{uuid}", svc.GetUser)
				r.Get("/", svc.ListUsers)
			})

			r.Group(func(r chi.Router) {
				r.Use(middleware.JWT)
				r.Use(middleware.ParseUUID, middleware.UUIDRights, middleware.IsTeacherRights)

				r.Patch("/{uuid}", svc.PatchUser)
				r.Delete("/{uuid}", svc.DeleteUser)

			})
			r.Post("/register", svc.Register)
			r.Post("/login", svc.Login)
		})
	})

	return r, nil
}
