package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/egot3/fathom/internal/carefulness"
	"github.com/google/uuid"
)

func ParseUUID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userUUID, err := uuid.Parse(r.PathValue("uuid"))
		if err != nil {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Couldn't parse requested uuid"})

			return
		}

		ctx := context.WithValue(r.Context(), "uuid", userUUID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
