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

		r.Route("/group", func(r chi.Router) {
			r.Get("/", svc.ListGroups)

			r.With(middleware.JWT, middleware.IsTeacherRights).Post("/", svc.PostGroup)

			r.Route("/{uuid}", func(r chi.Router) {
				r.Use(middleware.ParseUUID)
				r.Get("/", svc.GetGroup)

				r.Group(func(r chi.Router) {
					r.Use(middleware.JWT, middleware.IsTeacherRights)
					r.Patch("/", svc.PatchGroup)
					r.Delete("/", svc.DeleteGroup)
				})

				r.Route("/user", func(r chi.Router) {
					r.Post("/", svc.AppendUsers)
					r.Delete("/", svc.RemoveUsers)
				})
			})
		})

		r.Route("/quiz", func(r chi.Router) {
			r.Use(middleware.JWT, middleware.IsTeacherRights)

			r.Get("/", svc.ListQuizzes)
			r.Post("/", svc.PostQuiz)
			r.Post("/import", svc.ImportQuizBank)
			r.Get("/export", svc.ExportQuizBank)

			r.Route("/{uuid}", func(r chi.Router) {
				r.Use(middleware.ParseUUID)
				r.Patch("/", svc.PatchQuiz)
				r.Put("/", svc.PutQuiz)
				r.Get("/", svc.GetQuiz)
				r.Delete("/", svc.DeleteQuiz)
			})
		})

		r.Route("/test", func(r chi.Router) {
			r.Use(middleware.JWT)

			r.With(middleware.IsTeacherRights).Post("/", svc.PostTest)
			r.Get("/", svc.ListTests)

			r.Route("/running", func(r chi.Router) {
				r.Get("/", svc.GetQuizFromRunning)

				//protected
				r.Group(func(r chi.Router) {
					r.Use(middleware.IsTeacherRights)
					r.Post("/start", svc.StartTest)
					r.Post("/stop", svc.StopTest)
					r.Post("/pause", svc.PauseTest)
					r.Post("/resume", svc.ResumeTest)
					r.Post("/", svc.AddQuizzesToRunning)
					r.Delete("/", svc.RemoveQuizzesFromRunning)
					r.Post("/extend", svc.ExtendTest)
				})

			})

			r.Route("/{uuid}", func(r chi.Router) {
				r.Use(middleware.ParseUUID)

				r.Group(func(r chi.Router) {
					r.Use(middleware.IsTeacherRights)

					r.Delete("/", svc.DeleteTest)
					r.Patch("/", svc.PatchTest)
					r.Post("/quizzes", svc.AddQuizzes)
					r.Delete("/quizzes", svc.RemoveQuizzes)

				})

			})
		})
	})

	return r, nil
}
