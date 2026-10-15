package handler

import (
	"context"
	"encoding/csv"

	"github.com/egot3/fathom/internal/config"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

func CommitUserImport(ctx context.Context, reader *csv.Reader, svc chiService) error {
	return svc.commitUserImport(ctx, reader)
}

func SeviceStruct(i do.Injector) (chiService, error) {
	cf := do.MustInvoke[*config.Config](i)
	db := do.MustInvoke[*bun.DB](i)

	return chiService{
		cfg: cf,
		db:  db,
	}, nil
}
