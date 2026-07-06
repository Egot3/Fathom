package models

import (
	"context"
	"path/filepath"

	"github.com/egot3/fathom/internal/carefulness"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Quiz struct {
	bun.BaseModel `bun:"table:quizzes,alias:q"`

	UUID          uuid.UUID `bun:"uuid,pk"`
	Path          string    `bun:"path,unique"`
	Checksum      []byte    `bun:"checksum,notnull"`
	Score         int       `bun:"score,notnull"`
	CorrectAnswer string    `bun:"correct_answer,notnull"`
}

var _ bun.BeforeAppendModelHook = (*Quiz)(nil)

func (q *Quiz) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if !filepath.IsAbs(q.Path) {
			return carefulness.ErrAbsoluteRequired
		}
		if filepath.Ext(q.Path) != ".md" {
			return carefulness.PlainMarkdownRequired
		}
		q.UUID = uuid.New()
	}
	return nil
}
