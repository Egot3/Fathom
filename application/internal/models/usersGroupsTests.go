package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UserGroupsTests struct {
	bun.BaseModel `bun:"table:users_groups_tests,alias:ugt"`

	TestUUID uuid.UUID `bun:"test_uuid,pk"`
	Test     *Test     `bun:"rel:belongs-to,join:test_uuid=uuid"`

	GroupUUID uuid.UUID `bun:"group_uuid,pk"`
	Group     *Group    `bun:"rel:belongs-to,join:group_uuid=uuid"`

	UserUUID uuid.UUID `bun:"user_uuid,pk"`
	User     *User     `bun:"rel:belongs-to,join:user_uuid=uuid"`

	Score int `bun:"score"`
}
