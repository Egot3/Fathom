package jwtutils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func RemintToken(tokenString string) (string, error) {
	parser := jwt.NewParser(jwt.WithLeeway(30 * time.Minute))
	token, err := parser.ParseWithClaims(
		tokenString,
		&Claims{},

		func(t *jwt.Token) (any, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok { //not ok
				return nil, jwt.ErrSignatureInvalid
			}
			return secretKey, nil
		})
	if err != nil {
		return "", err
	}

	now := time.Now()
	expirationDate := now.Add(JWTTTL)

	claims, ok := token.Claims.(*Claims)
	newClaims := &Claims{
		UserID:    claims.UserID,
		IsTeacher: claims.IsTeacher,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expirationDate),
		},
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)

	if ok && token.Valid {
		return newToken.SignedString(secretKey)
	} else {
		return "", jwt.ErrTokenInvalidClaims
	}
}
