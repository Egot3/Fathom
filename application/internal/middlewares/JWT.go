package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	jwtutils "github.com/egot3/fathom/internal/JWTutils"
	"github.com/egot3/fathom/internal/carefulness"
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
		ctx := context.WithValue(rctx, "claims", *claims)

		next.ServeHTTP(w, r.WithContext(ctx))

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
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt_token",
			Value:    newToken,
			Path:     "/",
			Expires:  time.Now().Add(jwtutils.JWTTTL),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
		w.Header().Set("Session-Control", fmt.Sprintf("max-age: %d", int(math.Trunc(jwtutils.JWTTTL.Seconds()))))
	})
}
