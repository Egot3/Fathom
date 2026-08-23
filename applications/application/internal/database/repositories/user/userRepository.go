package user

import (
	"context"

	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	UpdateUser(ctx context.Context, patchedUser models.PatchUser) error
	IsTeacher(ctx context.Context, uuid uuid.UUID) (bool, error)
	DeleteUser(ctx context.Context, uuid uuid.UUID) error
	Exists(ctx context.Context, uuid uuid.UUID) (bool, error)
	Login(ctx context.Context, nickname string, passwordHash []byte) (models.User, error)
	Register(ctx context.Context, name string, passwordHash []byte) (models.User, error)
	User(ctx context.Context, uuid uuid.UUID) (models.User, error)
	List(ctx context.Context, page, size int) ([]models.User, int, error)
}
