package jwtutils

import (
	"context"

	"github.com/google/uuid"
)

type uuidKey struct{}

func UUIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	UUID, ok := ctx.Value(uuidKey{}).(uuid.UUID)
	return UUID, ok
}

func WithUUID(ctx context.Context, UUID uuid.UUID) context.Context {
	return context.WithValue(ctx, uuidKey{}, UUID)
}
