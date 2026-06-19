package jwtutils

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct { //Метаморфоз https://purpleschool.ru/knowledge-base/article/work-with-jwt
	UserID    uuid.UUID `json:"user_uuid"`
	IsTeacher bool      `json:"is_teacher"`
	jwt.RegisteredClaims
}
