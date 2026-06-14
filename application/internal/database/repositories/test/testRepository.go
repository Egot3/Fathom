package test

import (
	"context"

	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
)

type TestRepository interface {
	CreateTest(ctx context.Context, name string) error
	BundleQuizzesToTest(ctx context.Context, testUUID uuid.UUID, pathes []string) error
	PruneQuizzesFromTest(ctx context.Context, testUUID uuid.UUID, pathes []string) error
	Test(ctx context.Context, UUID uuid.UUID) (*models.Test, error)
	DeleteTest(ctx context.Context, UUID uuid.UUID) error
	UpdateTest(ctx context.Context, UUID uuid.UUID, name string) error
}
