package starters

import (
	"github.com/egot3/fathom/internal/models"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

func registerAll(i do.Injector) {
	db := do.MustInvoke[*bun.DB](i)

	models.RegisterAll(db)
}
