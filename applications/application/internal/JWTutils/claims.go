package jwtutils

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID    uuid.UUID `json:"user_uuid"`
	IsTeacher bool      `json:"is_teacher"`
	jwt.RegisteredClaims
}
