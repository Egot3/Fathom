package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Test struct {
	bun.BaseModel `bun:"table:tests,alias:t"`

	UUID uuid.UUID `bun:"uuid,pk" json:"uuid"`
	Name string    `bun:"name" json:"name"`

	CreatedAt time.Time `bun:"created_at,notnull,default:now()" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:now()" json:"updated_at"`

	Quizzes []Quiz `bun:"m2m:tests_quizzes,join:Test=Quiz" json:"quizzes"`
}

var _ bun.BeforeAppendModelHook = (*Test)(nil)

func (t *Test) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		var err error
		t.UUID, err = uuid.NewV7() //compensating sqlite
		if err != nil {
			return err
		}
		t.CreatedAt = time.Now()
		t.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		t.UpdatedAt = time.Now()
	}
	return nil
}
