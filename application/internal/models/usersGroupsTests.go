package models

import (
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

	Score int `bun:"score" json:"score"`
}
