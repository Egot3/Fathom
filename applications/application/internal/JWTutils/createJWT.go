package jwtutils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const JWTTTL = time.Hour

func GenerateToken(userID uuid.UUID, isTeacher bool) (string, error) {
	now := time.Now()
	expirationDate := now.Add(JWTTTL)

	claims := &Claims{
		UserID:    userID,
		IsTeacher: isTeacher,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expirationDate),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}
