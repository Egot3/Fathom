package models

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Group struct {
	bun.BaseModel `bun:"table:groups,alias:g"`

	UUID uuid.UUID `bun:"uuid,pk"`
	Name string    `bun:"name,unique"`
}

var _ bun.BeforeAppendModelHook = (*Group)(nil)

func (g *Group) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		g.UUID = uuid.New() //compensating sqlite
	}
	return nil
}
