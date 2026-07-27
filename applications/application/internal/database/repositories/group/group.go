package group

import (
	"context"
	"database/sql"
	"errors"

	"github.com/egot3/fathom/internal/carefulness"
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

func (r *bunGroupRepository) NewGroup(ctx context.Context, name string) (*models.Group, error) {
	var group = models.Group{Name: name}
	err := r.db.NewInsert().Model(&group).Ignore().Returning("uuid").Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, carefulness.Conflict{Conflictor: "group name"}
		}
		return nil, err
	}
	return &group, nil
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
	e, err := r.db.NewSelect().Model((*models.Group)(nil)).
		Where("uuid <> ?", uuid).Where("name = ?", name).
		Exists(ctx)
	if err != nil {
		return err
	}
	if e {
		return carefulness.Conflict{Conflictor: "name"}
	}

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
		res, err := tx.NewInsert().Ignore().Model(&groupUsers).Exec(ctx)
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
		if int(c) < len(userUUIDs) {
			return carefulness.PartialSuccess{
				Target:      "users",
				ActualCount: int(c),
				WantCount:   len(userUUIDs),
			}
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

func (r *bunGroupRepository) ListGroups(ctx context.Context, page, size int) ([]models.Group, int, error) {
	var groups []models.Group
	total, err := r.db.NewSelect().Model(&groups).
		OrderBy("name", bun.OrderAsc).
		Offset(size * page).Limit(size).
		ScanAndCount(ctx)
	if err != nil {

		return nil, 0, err
	}

	return groups, total, nil
}

func (r *bunGroupRepository) GroupsExist(ctx context.Context, groupUUIDs uuid.UUIDs) (bool, error) {
	count, err := r.db.NewSelect().
		Model((*models.Group)(nil)).
		Where("uuid IN (?)", bun.List(groupUUIDs)).
		Group("uuid").
		Having("COUNT(DISTINCT uuid) = ?", len(groupUUIDs)).
		Count(ctx)
	if err != nil {
		return false, err
	}

	allExist := count == len(groupUUIDs)
	if !allExist {
		return false, nil
	}

	return true, nil
}

func (r *bunGroupRepository) IsInAny(ctx context.Context, groupUUIDs uuid.UUIDs, userUUID uuid.UUID) (bool, error) {
	is, err := r.db.NewSelect().
		Model((*models.GroupsUsers)(nil)).
		Where("group_uuid IN (?)", bun.List(groupUUIDs)).
		Where("user_uuid = ?", userUUID).Exists(ctx)
	if err != nil {
		return false, err
	}

	return is, nil
}
