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
	err := r.db.NewSelect().Model(&group).WherePK().Relation("User").Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (r *bunGroupRepository) DeleteGroup(ctx context.Context, uuid uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&models.Group{UUID: uuid}).WherePK().Exec(ctx)
	return err
}

func (r *bunGroupRepository) UpdateGroup(ctx context.Context, uuid uuid.UUID, name string) error {
	_, err := r.db.NewUpdate().Model(&models.Group{UUID: uuid, Name: name}).WherePK().Exec(ctx)
	return err
}

func (r *bunGroupRepository) AppendUsers(ctx context.Context, groupUUID uuid.UUID, userUUIDs uuid.UUIDs) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, userUUID := range userUUIDs {
			if _, err := tx.NewInsert().Model(&models.GroupsUsers{GroupUUID: groupUUID, UserUUID: userUUID}).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *bunGroupRepository) RemoveUsers(ctx context.Context, groupUUID uuid.UUID, userUUIDs uuid.UUIDs) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, userUUID := range userUUIDs {
			if _, err := tx.NewDelete().Model(&models.GroupsUsers{GroupUUID: groupUUID, UserUUID: userUUID}).WherePK().Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}
