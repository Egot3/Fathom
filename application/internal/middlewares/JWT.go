package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"time"

	jwtutils "github.com/egot3/fathom/internal/JWTutils"
	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/logging"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.LoggerFromContext(r.Context()).With(slog.String("layer", "middleware"))

		authorization, err := r.Cookie("jwt_token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		logger.Debug("passed token cookie getting", slog.String("token", authorization.Value))

		claims, err := jwtutils.ValidateToken(authorization.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			cookie := &http.Cookie{
				Name:     "jwt_token",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(w, cookie)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Bad token"})
			return
		}
		logger.Debug("token validated")

		ctx := r.Context()
		ctx = context.WithValue(ctx, "claims", *claims)

		newToken, err := jwtutils.RemintToken(authorization.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			cookie := &http.Cookie{
				Name:     "jwt_token",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(w, cookie)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Bad token"})
			return
		}
		logger.Debug("token reminted", slog.String("newToken", newToken))

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt_token",
			Value:    newToken,
			Path:     "/",
			Expires:  time.Now().Add(jwtutils.JWTTTL),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		mAge := int(math.Trunc(jwtutils.JWTTTL.Seconds()))
		logger.Debug("sent token", slog.Int("maxAge", mAge))

		w.Header().Set("Session-Control", fmt.Sprintf("max-age=%d", mAge))
		logger.Debug("Sent Session-Control header")

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
