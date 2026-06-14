package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	UUID         uuid.UUID `bun:"uuid,pk"`
	Nickname     string    `bun:"nickname,notnull,unique"`
	PasswordHash []byte    `bun:"password_hash,notnull"`
	IsTeacher    bool      `bun:"is_teacher,default:false"`

	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`
	CreatedAt time.Time  `bun:"created_at"`
	UpdatedAt time.Time  `bun:"updated_at"`
}

var _ bun.BeforeAppendModelHook = (*User)(nil)

func (u *User) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		u.UUID = uuid.New() //compensating sqlite
		u.UpdatedAt = time.Now()
		u.DeletedAt = nil
		u.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		u.UpdatedAt = time.Now()
	}
	return nil
}

type PatchUser struct {
	UUID uuid.UUID

	Nickname     *string
	PasswordHash []byte
	IsTeacher    *bool
}
