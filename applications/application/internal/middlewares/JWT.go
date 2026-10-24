package middlewares

import (
	"context"
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

		authorization, err := r.Cookie("jwt_token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := jwtutils.ValidateToken(authorization.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			http.SetCookie(w, httputils.ResetCookie)
			json.NewEncoder(w).Encode(&carefulness.JSONError{Err: "Bad token"})
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "claims", *claims)

		cookie, err := httputils.NewAuthCookie(claims.UserID, claims.IsTeacher)
		http.SetCookie(w, cookie)
		if err != nil {
			logging.LoggerFromContext(r.Context()).Error("unable to create login token", slog.String("Error", err.Error()))
			(&carefulness.JSONError{Err: "Unable to create login token", Status: http.StatusInternalServerError}).Encode(w)
			return
		}

		w.Header().Set("Session-Control", fmt.Sprintf("max-age=%d", cookie.MaxAge))

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
