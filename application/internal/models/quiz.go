package models

import (
	"github.com/uptrace/bun"
)

type Quiz struct {
	bun.BaseModel `bun:"table:quizes,alias:q"`

	Path string `bun:"path,unique"`
}
