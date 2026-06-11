package quiz

import (
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunQuizRepository struct {
	db *bun.DB
}

func NewQuizRepository(i do.Injector)
