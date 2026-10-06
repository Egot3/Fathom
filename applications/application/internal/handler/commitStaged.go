package handler

import (
	"context"
	"os"
	"path/filepath"

	"github.com/egot3/fathom/internal/importutils"
	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
)

func (c *chiService) commitTestImport(ctx context.Context, staged importutils.StagedTest) error {
	quizTestPairs := lo.Map(staged.QuizUUIDs, func(UUID uuid.UUID, i int) models.TestsQuizzes {
		return models.TestsQuizzes{
			TestUUID: staged.Test.UUID,
			QuizUUID: UUID,
			Position: i,
		}
	})
	return c.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(&staged.Test).Exec(ctx); err != nil {
			return err
		}

		_, err := tx.NewInsert().Model(&quizTestPairs).Exec(ctx)
		return err
	})
}

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
