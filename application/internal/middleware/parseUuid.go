package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/logging"
	"github.com/google/uuid"
)

func ParseUUID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.LoggerFromContext(r.Context())

		UUID, err := uuid.Parse(r.PathValue("uuid"))
		if err != nil {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Couldn't parse requested uuid"})

			return
		}

		logger.Info("parsed uuid in middelware",
			slog.String("uuid", UUID.String()))
		ctx := context.WithValue(r.Context(), "uuid", UUID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
