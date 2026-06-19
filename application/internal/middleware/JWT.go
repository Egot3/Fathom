package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	jwtutils "github.com/egot3/fathom/internal/JWTutils"
	"github.com/egot3/fathom/internal/carefulness"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := jwtutils.ValidateToken(token)
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
		}

		rctx := r.Context()
		ctx := context.WithValue(rctx, "claims", &claims)

		next.ServeHTTP(w, r.WithContext(ctx))

		newToken, err := jwtutils.RemintToken(token)
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
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt_token",
			Value:    newToken,
			Path:     "/",
			Expires:  time.Now().Add(jwtutils.JWTTTL),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
	})
}
