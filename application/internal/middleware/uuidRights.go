package middleware

import (
	"encoding/json"
	"net/http"

	jwtutils "github.com/egot3/fathom/internal/JWTutils"
	"github.com/egot3/fathom/internal/carefulness"
	"github.com/google/uuid"
)

// requires JWT
func Rights(next http.Handler) http.Handler { // did you know you have rights?
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := (r.Context().Value("claims")).(jwtutils.Claims)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Unable to retrieve jwt's claims"})
			return
		}

		uuid, ok := (r.Context().Value("uuid")).(uuid.UUID)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Unable to retrieve uuid"})
			return
		}

		if claims.UserID != uuid && !claims.IsTeacher {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Not enough permissions for operation"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
