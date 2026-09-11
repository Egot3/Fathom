package testutils

import (
	"math/rand/v2"
	"strings"
)

// safe base64
var charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

func GenerateRandomString(size int) string {
	var sb strings.Builder
	sb.Grow(size)

	for i := 0; i < size; i++ {
		randomIndex := rand.IntN(len(charset))
		sb.WriteByte(charset[randomIndex])
	}

	return sb.String()
}
