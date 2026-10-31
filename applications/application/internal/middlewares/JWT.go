package middlewares

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	jwtutils "github.com/egot3/fathom/internal/JWTutils"
	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/httputils"
	"github.com/egot3/fathom/internal/logging"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.LoggerFromContext(ctx)

		authorization, err := r.Cookie("jwt")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := jwtutils.ValidateToken(authorization.Value)
		if err != nil {
			logger.Warn("got bad token")
			w.WriteHeader(http.StatusUnauthorized)
			http.SetCookie(w, httputils.ResetCookie)
			json.NewEncoder(w).Encode(&carefulness.JSONError{Err: "Bad token"})
			return
		}

		ctx = jwtutils.WithClaims(ctx, *claims)
		ctx = logging.WithLogger(ctx, logger.With(slog.String("user_uuid", claims.UserID.String())))

		cookie, err := httputils.NewAuthCookie(claims.UserID, claims.IsTeacher)
		http.SetCookie(w, cookie)
		if err != nil {
			logger.Error("unable to create login token", slog.String("Error", err.Error()))
			(&carefulness.JSONError{Err: "Unable to create login token", Status: http.StatusInternalServerError}).Encode(w)
			return
		}

		w.Header().Set("Session-Control", fmt.Sprintf("max-age=%d", cookie.MaxAge))

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
