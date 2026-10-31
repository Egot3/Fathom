package middlewares

import (
	"encoding/json"
	"log/slog"
	"net/http"

	jwtutils "github.com/egot3/fathom/internal/JWTutils"
	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/logging"
	"github.com/go-chi/chi/v5"
)

// requires JWT
// requires parseUUID
func UUIDRights(next http.Handler) http.Handler { // did you know you have rights?
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

		uuid, ok := jwtutils.UUIDFromContext(ctx)
		if !ok {
			logger.Error("Failed to retrieve uuid")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&carefulness.JSONError{Err: "Unable to retrieve uuid"})
			return
		}

		if claims.UserID != uuid {
			logger.Warn("No rights for operation",
				slog.String("requestorUUID", claims.ID),
				slog.String("requestedUUID", uuid.String()),
			)
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(&carefulness.JSONError{Err: "Not enough permissions for operation"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func UserRights(next http.Handler) http.Handler {
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

		if claims.UserID.String() != chi.URLParam(r, "user_uuid") {
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

func IsUserRights(r *http.Request) bool {
	claims, ok := jwtutils.ClaimsFromContext(r.Context())
	if !ok {
		return false
	}

	if claims.UserID.String() != chi.URLParam(r, "user_uuid") {
		return false
	}

	return true
}
