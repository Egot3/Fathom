package jwtutils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var JWTTTL = 3600 * time.Second

func GenerateToken(userID uuid.UUID, isTeacher bool) (string, error) {
	now := time.Now()
	expirationDate := now.Add(JWTTTL)

	claims := &Claims{
		UserID:    userID, //неожидано
		IsTeacher: isTeacher,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expirationDate),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}
