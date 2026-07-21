package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UserGroupsTests struct {
	bun.BaseModel `bun:"table:users_groups_tests,alias:ugt"`

	TestUUID uuid.UUID `bun:"test_uuid,pk" json:"test_uuid"`
	Test     *Test     `bun:"rel:belongs-to,join:test_uuid=uuid"`

	GroupUUID uuid.UUID `bun:"group_uuid,pk" json:"group_uuid"`
	Group     *Group    `bun:"rel:belongs-to,join:group_uuid=uuid"`

	UserUUID uuid.UUID `bun:"user_uuid,pk" json:"user_uuid"`
	User     *User     `bun:"rel:belongs-to,join:user_uuid=uuid"`

	Score       float64   `bun:"score" json:"score"`
	FinalizedAt time.Time `bun:"finalized_at" json:"finalized_at"`
}

var _ bun.BeforeAppendModelHook = (*UserGroupsTests)(nil)

func (ugt *UserGroupsTests) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		ugt.FinalizedAt = time.Now()
	}
	return nil
}
