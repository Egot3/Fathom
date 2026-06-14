package repositories

import (
	"context"

	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunUserRepository struct {
	db *bun.DB
}

func NewUserRepository(i do.Injector) (UserRepository, error) {
	db := do.MustInvoke[*bun.DB](i)

	return &bunUserRepository{db: db}, nil
}

func (r *bunUserRepository) Register(ctx context.Context, name string, passwordHash []byte) error {
	_, err := r.db.NewInsert().Model(&models.User{Nickname: name, PasswordHash: passwordHash}).Exec(ctx)
	return err
}

func (r *bunUserRepository) Login(ctx context.Context, uuid uuid.UUID, passwordHash []byte) (success bool, err error) {
	success, err = r.db.NewSelect().Model(&models.User{UUID: uuid}).WherePK().Where("password_hash = ?", passwordHash).Exists(ctx)
	return success, err // on err success is false(by bun's code)
}

func (r *bunUserRepository) Exists(ctx context.Context, uuid uuid.UUID) (bool, error) {
	exists, err := r.db.NewSelect().Model(&models.User{UUID: uuid}).WherePK().Exists(ctx)
	return exists, err
}

func (r *bunUserRepository) DeleteUser(ctx context.Context, uuid uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&models.User{UUID: uuid}).WherePK().Exec(ctx)
	return err
}

func (r *bunUserRepository) IsTeacher(ctx context.Context, uuid uuid.UUID) (bool, error) {
	is, err := r.db.NewSelect().Model(&models.User{UUID: uuid}).WherePK().Where("is_teacher = ?", true).Exists(ctx)
	return is, err
}

func (r *bunUserRepository) User(ctx context.Context, uuid uuid.UUID) (*models.User, error) {
	user := models.User{UUID: uuid}
	err := r.db.NewSelect().Model(&user).WherePK().Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &models.User{
		Nickname:     user.Nickname,
		UUID:         user.UUID,
		IsTeacher:    user.IsTeacher,
		PasswordHash: nil, //safety
	}, nil
}

func (r *bunUserRepository) UpdateUser(ctx context.Context, patchedUser models.PatchUser) error {
	query := r.db.NewUpdate().Model(&models.User{UUID: patchedUser.UUID}).WherePK()
	if patchedUser.Nickname != nil {
		query = query.Set("nickname = ?", patchedUser.Nickname)
	}
	if patchedUser.PasswordHash != nil {
		query = query.Set("password_hash = ?", patchedUser.PasswordHash)
	}
	if patchedUser.IsTeacher != nil {
		query = query.Set("is_teacher = ?", patchedUser.IsTeacher)
	}

	_, err := query.Exec(ctx)

	return err
}
