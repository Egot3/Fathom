package group

import (
	"context"

	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
)

type GroupRepository interface {
	NewGroup(ctx context.Context, name string) error
	DeleteGroup(ctx context.Context, uuid uuid.UUID) error
	Group(ctx context.Context, uuid uuid.UUID) (*models.Group, error)
	UpdateGroup(ctx context.Context, uuid uuid.UUID, name string) error
	AppendUsers(ctx context.Context, groupUUID uuid.UUID, userUUIDs uuid.UUIDs) error
	RemoveUsers(ctx context.Context, groupUUID uuid.UUID, userUUIDs uuid.UUIDs) error
	IsInGroup(ctx context.Context, groupUUID, userUUID uuid.UUID) (bool, error)
	ListGroups(ctx context.Context, page, size int) ([]models.Group, int, error)
}
