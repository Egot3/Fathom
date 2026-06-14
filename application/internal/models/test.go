package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Test struct {
	bun.BaseModel `bun:"table:tests,alias:t"`

	UUID uuid.UUID `bun:"uuid,pk"`
	Name string    `bun:"name"`

	CreatedAt time.Time `bun:"created_at,notnull,default:now()"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:now()"`

	Quizzes []Quiz `bun:"m2m:tests_quizzes,join:Test=Quiz"`
}

var _ bun.BeforeAppendModelHook = (*Test)(nil)

func (t *Test) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		t.UUID = uuid.New() //compensating sqlite
		t.CreatedAt = time.Now()
		t.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		t.UpdatedAt = time.Now()
	}
	return nil
}
