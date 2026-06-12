package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type GroupsUsers struct {
	bun.BaseModel `bun:"table:groups_users,alias:gu"`

	GroupUUID uuid.UUID `bun:"group_uuid,pk"`
	Group     *Group    `bun:"rel:belongs-to,join:group_uuid=uuid"`

	UserUUID uuid.UUID `bun:"user_uuid,pk"`
	User     *User     `bun:"rel:belongs-to,join:user_uuid=uuid"`
}
