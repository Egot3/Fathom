package models

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Test struct {
	bun.BaseModel `bun:"table:tests,alias:t"`

	UUID     uuid.UUID `bun:"uuid,pk"`
	Name     string    `bun:"name"`
	MaxScore string    `bun:"name"` //computed
}

var _ bun.BeforeAppendModelHook = (*Test)(nil)

func (t *Test) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		t.UUID = uuid.New() //compensating sqlite
	}
	return nil
}
