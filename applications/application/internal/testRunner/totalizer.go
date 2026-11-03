package testrunner

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Totalizer interface {
	TotalizeAll(ctx context.Context, testUUID uuid.UUID, groupUUIDs uuid.UUIDs, startedAt time.Time) error
}

type nopTotalizer struct{}

func (nopTotalizer) TotalizeAll(ctx context.Context, testUUID uuid.UUID, groupUUIDs uuid.UUIDs, startedAt time.Time) error {
	return nil // does nothing -- profit
}
