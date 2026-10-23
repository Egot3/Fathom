package httputils

import (
	"net/http"
	"time"

	jwtutils "github.com/egot3/fathom/internal/JWTutils"
	"github.com/google/uuid"
)

var ResetCookie = &http.Cookie{
	Name:        "jwt_token",
	Value:       "",
	Path:        "/",
	Expires:     time.Unix(0, 0),
	HttpOnly:    true,
	SameSite:    http.SameSiteNoneMode,
	Secure:      true,
	Partitioned: true,
}

func NewAuthCookie(userUUID uuid.UUID, isTeacher bool) (*http.Cookie, error) {
	newToken, err := jwtutils.GenerateToken(userUUID, isTeacher)
	if err != nil {
		return ResetCookie, err
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
	if err := cookie.Valid(); err != nil {
		return ResetCookie, err
	}
	return cookie, nil
}
