package models

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/egot3/fathom/internal/carefulness"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Quiz struct {
	bun.BaseModel `bun:"table:quizzes,alias:q"`

	UUID          uuid.UUID   `bun:"uuid,pk" json:"uuid"`
	Path          string      `bun:"path,unique" json:"path"`
	Checksum      FixedBytes8 `bun:"checksum,notnull"`
	Score         int         `bun:"score,notnull" json:"score"`
	CorrectAnswer string      `bun:"correct_answer,notnull" json:"correct_answer"`
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
		var err error
		q.UUID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return nil
}

type FixedBytes8 [8]byte

func (f FixedBytes8) Value() (driver.Value, error) {
	return bzCopy(f[:]), nil
}

func (f *FixedBytes8) Scan(value any) error {
	switch b := value.(type) {
	case []byte:
		if len(b) != 8 {
			return fmt.Errorf("invalid length for [8]byte: %d", len(b))
		}
		copy(f[:], b)
		return nil
	case int64:
		for i := 7; i >= 0; i-- {
			f[i] = byte(b >> uint(8*(7-i)))
		}
		return nil
	case FixedBytes8:
		b = *f
		return nil
	}
	return errors.New("failed to scan [8]byte: incompatible type")
}

func bzCopy(b []byte) []byte {
	out := make([]byte, len(b))
	copy(out, b)
	return out
}
