package jwtutils

import (
	"context"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

type jwtKey struct{}

type Claims struct {
	UserID    uuid.UUID `json:"user_uuid"`
	IsTeacher bool      `json:"is_teacher"`
	jwt.RegisteredClaims
}

func WithClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, jwtKey{}, claims)
}

func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	cl, ok := ctx.Value(jwtKey{}).(Claims)
	return cl, ok
}
