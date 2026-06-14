package models

import (
	"github.com/uptrace/bun"
)

type Quiz struct {
	bun.BaseModel `bun:"table:quizzes,alias:q"`

	Path     string `bun:"path,pk"`
	Checksum []byte `bun:"checksum,notnull"`
	Score    int    `bun:"score,notnull"`
}
