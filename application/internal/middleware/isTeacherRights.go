package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"

	jwtutils "github.com/egot3/fathom/internal/JWTutils"
	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/logging"
)

// requires JWT
func IsTeacherRights(next http.Handler) http.Handler { // did you know you don't have rights?
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.LoggerFromContext(r.Context())
		logger = logger.With(slog.String("layer", "middleware"))

		claims, ok := (r.Context().Value("claims")).(jwtutils.Claims)
		if !ok {
			logger.Error("Failed to retrieve jwt claims")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Unable to retrieve jwt's claims"})
			return
		}

		if !claims.IsTeacher {
			logger.Warn("No rights for operation",
				slog.String("requestorUUID", claims.ID),
				slog.Bool("IsTeacher", claims.IsTeacher),
			)
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Not enough permissions for operation"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
