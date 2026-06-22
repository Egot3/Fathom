package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type bunUserRepository struct {
	db *bun.DB
}

func NewUserRepository(i do.Injector) (UserRepository, error) {
	db := do.MustInvoke[*bun.DB](i)

	return &bunUserRepository{db: db}, nil
}

func (r *bunUserRepository) Register(ctx context.Context, name string, passwordHash []byte) (*models.User, error) {
	bakedUser := models.User{Nickname: name, PasswordHash: passwordHash}
	err := r.db.NewInsert().Ignore().Model(&bakedUser).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, carefulness.Conflict{Conflictor: name}
		}
		return nil, err
	}

	return &models.User{
		UUID:      bakedUser.UUID,
		Nickname:  bakedUser.Nickname,
		CreatedAt: bakedUser.CreatedAt,
		UpdatedAt: bakedUser.UpdatedAt,
		DeletedAt: bakedUser.DeletedAt,
		IsTeacher: bakedUser.IsTeacher,
	}, nil
}

func (r *bunUserRepository) Login(ctx context.Context, nickname string, password []byte) (*models.User, error) {
	var user models.User
	err := r.db.NewSelect().Model(&user).
		Where("nickname = ?", nickname).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, password)
	if err != nil {
		return nil, sql.ErrNoRows // will be sql.ErrNoRows
	}
	return &models.User{
		UUID:      user.UUID,
		Nickname:  user.Nickname,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
		IsTeacher: user.IsTeacher,
	}, nil
}

func (r *bunUserRepository) Exists(ctx context.Context, uuid uuid.UUID) (bool, error) {
	exists, err := r.db.NewSelect().Model(&models.User{UUID: uuid}).WherePK().Exists(ctx)
	return exists, err
}

func (r *bunUserRepository) List(ctx context.Context, page, size int) ([]models.User, int, error) {
	var users []models.User
	total, err := r.db.NewSelect().
		Model(&users).
		OrderBy("created_at", bun.OrderAsc).
		ScanAndCount(ctx)
	if err != nil {
		return nil, total, err
	}

	return users, total, nil
}

func (r *bunUserRepository) DeleteUser(ctx context.Context, uuid uuid.UUID) error {
	res, err := r.db.NewDelete().Model(&models.User{UUID: uuid}).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	c, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if c == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *bunUserRepository) IsTeacher(ctx context.Context, uuid uuid.UUID) (bool, error) {
	is, err := r.db.NewSelect().Model(&models.User{UUID: uuid}).WherePK().Where("is_teacher = ?", true).Exists(ctx)
	return is, err
}

func (r *bunUserRepository) User(ctx context.Context, uuid uuid.UUID) (*models.User, error) {
	user := models.User{UUID: uuid}
	err := r.db.NewSelect().Model(&user).WherePK().WhereAllWithDeleted().Scan(ctx)
	if err != nil {
		return nil, err
	}

	if user.DeletedAt != nil && time.Since(*user.DeletedAt) >= 0 {
		return nil, carefulness.ErrGone
	}

	return &models.User{
		Nickname:     user.Nickname,
		UUID:         user.UUID,
		IsTeacher:    user.IsTeacher,
		PasswordHash: nil, //safety
	}, nil
}

func (r *bunUserRepository) UpdateUser(ctx context.Context, patchedUser models.PatchUser) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		query := tx.NewUpdate().Model((*models.User)(nil)).
			Where("uuid = ?", patchedUser.UUID)
		if patchedUser.Nickname != nil {
			e, err := tx.NewSelect().Model((*models.User)(nil)).
				Where("uuid <> ?", patchedUser.UUID).
				Where("nickname = ?", patchedUser.Nickname).
				Exists(ctx)
			if err != nil {
				return fmt.Errorf("Error while checking for unique: %w", err)
			}
			if e {
				return carefulness.Conflict{Conflictor: "nickname"}
			}
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
	})
}
