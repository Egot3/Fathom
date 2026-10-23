package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwtutils "github.com/egot3/fathom/internal/JWTutils"
	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/httputils"
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

		newToken, err := jwtutils.RemintToken(authorization.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			http.SetCookie(w, httputils.ResetCookie)
			json.NewEncoder(w).Encode(&carefulness.JSONError{Err: "Bad token"})
			return
		}

		cookie := &http.Cookie{
			Name:        "jwt_token",
			Value:       newToken,
			Path:        "/",
			Expires:     time.Now().Add(jwtutils.JWTTTL),
			HttpOnly:    true,
			SameSite:    http.SameSiteNoneMode,
			Secure:      true,
			Partitioned: true,
		}
		http.SetCookie(w, cookie)

		w.Header().Set("Session-Control", fmt.Sprintf("max-age=%d", cookie.MaxAge))

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
