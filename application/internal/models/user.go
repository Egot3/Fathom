package models

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	UUID         uuid.UUID `bun:"uuid,pk"`
	Nickname     string    `bun:"nickname,notnull,unique"`
	PasswordHash []byte    `bun:"password_hash,notnull"`
	IsTeacher    bool      `bun:"is_teacher,default:false"`
}

var _ bun.BeforeAppendModelHook = (*User)(nil)

func (u *User) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		u.UUID = uuid.New() //compensating sqlite
	}
	return nil
}
