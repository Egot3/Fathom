package handler

import (
	"context"
	"os"
	"path/filepath"

	"github.com/egot3/fathom/internal/importutils"
	"github.com/egot3/fathom/internal/models"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
)

func (c *chiService) commitImport(ctx context.Context, staged []importutils.StagedQuiz) error {
	var placed []string

	err := c.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		quizzes := lo.Map(staged, func(s importutils.StagedQuiz, _ int) models.Quiz { return s.Quiz })

		if _, err := tx.NewInsert().Model(&quizzes).Exec(ctx); err != nil {
			return err
		}

		for _, s := range staged {
			if err := os.MkdirAll(filepath.Dir(s.Quiz.Path), 0o755); err != nil {
				return err
			}
			if err := os.Link(s.StagedPath, s.Quiz.Path); err != nil {
				return err
			}
			placed = append(placed, s.Quiz.Path)
		}
		return nil
	})
	if err != nil {
		for _, p := range placed {
			_ = os.Remove(p)
		}
		return err
	}
	return nil
}
