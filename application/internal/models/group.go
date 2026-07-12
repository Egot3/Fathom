package models

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Group struct {
	bun.BaseModel `bun:"table:groups,alias:g"`

	UUID uuid.UUID `bun:"uuid,pk" json:"uuid"`
	Name string    `bun:"name,unique" json:"name"`

	Users []User `bun:"m2m:groups_users,join:Group=User" json:"pupils"`
}

var _ bun.BeforeAppendModelHook = (*Group)(nil)

func (g *Group) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		var err error
		g.UUID, err = uuid.NewV7() //compensating sqlite
		if err != nil {
			return err
		}
	}
	return nil
}
