package middlewares

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
		ctx := r.Context()
		logger := logging.LoggerFromContext(ctx)
		logger = logger.With(slog.String("layer", "middleware"))

		claims, ok := jwtutils.ClaimsFromContext(ctx)
		if !ok {
			logger.Error("Failed to retrieve jwt claims")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&carefulness.JSONError{Err: "Unable to retrieve jwt's claims"})
			return
		}

		if !claims.IsTeacher {
			logger.Warn("No rights for operation",
				slog.String("requestorUUID", claims.ID),
				slog.Bool("IsTeacher", claims.IsTeacher),
			)
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(&carefulness.JSONError{Err: "Not enough permissions for operation"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

// still requires JWT
func IsTeacherCondition(r *http.Request) bool {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)
	logger = logger.With(slog.String("layer", "condition"))

	claims, ok := jwtutils.ClaimsFromContext(ctx)
	if !ok {
		logger.Error("Failed to retrieve jwt claims")
		return false
	}

	return claims.IsTeacher
}
