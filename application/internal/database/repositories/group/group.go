package group

import (
	"context"
	"database/sql"

	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunGroupRepository struct {
	db *bun.DB
}

func NewGroupRepository(i do.Injector) (GroupRepository, error) {
	db := do.MustInvoke[*bun.DB](i)
	return &bunGroupRepository{db: db}, nil
}

func (r *bunGroupRepository) NewGroup(ctx context.Context, name string) error {
	_, err := r.db.NewInsert().Model(&models.Group{Name: name}).Exec(ctx)
	return err
}

func (r *bunGroupRepository) Group(ctx context.Context, uuid uuid.UUID) (*models.Group, error) {
	group := models.Group{UUID: uuid}
	err := r.db.NewSelect().Model(&group).WherePK().Relation("Users").Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (r *bunGroupRepository) DeleteGroup(ctx context.Context, uuid uuid.UUID) error {
	res, err := r.db.NewDelete().Model(&models.Group{UUID: uuid}).WherePK().Exec(ctx)
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

func (r *bunGroupRepository) UpdateGroup(ctx context.Context, uuid uuid.UUID, name string) error {
	res, err := r.db.NewUpdate().Model(&models.Group{UUID: uuid, Name: name}).WherePK().Exec(ctx)
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

// updated ver
func (r *bunGroupRepository) AppendUsers(ctx context.Context, groupUUID uuid.UUID, userUUIDs uuid.UUIDs) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		var groupUsers []models.GroupsUsers = make([]models.GroupsUsers, len(userUUIDs)) // updated after test
		for i, userUUID := range userUUIDs {
			groupUsers[i] = models.GroupsUsers{UserUUID: userUUID, GroupUUID: groupUUID}
		}
		if _, err := tx.NewInsert().Model(&groupUsers).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}

// updated after append
func (r *bunGroupRepository) RemoveUsers(ctx context.Context, groupUUID uuid.UUID, userUUIDs uuid.UUIDs) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		var groupUsers []models.GroupsUsers = make([]models.GroupsUsers, len(userUUIDs))
		for i, userUUID := range userUUIDs {
			groupUsers[i] = models.GroupsUsers{UserUUID: userUUID, GroupUUID: groupUUID}
		}

		res, err := tx.NewDelete().Model(&groupUsers).WherePK().Exec(ctx)
		if err != nil {
			return err
		}

		c, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if int(c) < len(userUUIDs) {
			return sql.ErrNoRows
		}
		return nil
	})
}

func (r *bunGroupRepository) IsInGroup(ctx context.Context, groupUUID, userUUID uuid.UUID) (bool, error) {
	return r.db.NewSelect().Model(&models.GroupsUsers{GroupUUID: groupUUID, UserUUID: userUUID}).
		WherePK().Exists(ctx)
}
