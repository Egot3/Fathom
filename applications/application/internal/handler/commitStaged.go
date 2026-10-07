package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/importutils"
	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
)

type UUIDPair struct {
	ExistingUUID uuid.UUID `bun:"existing_uuid"`
	AddedUUID    uuid.UUID `bun:"added_uuid"`
}

func (c *chiService) commitTestImport(
	ctx context.Context,
	quizzes []importutils.StagedQuiz,
	st importutils.StagedTest,
) error {
	return c.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		resolved := make(map[string]uuid.UUID, len(quizzes))

		for i := range quizzes {
			q := &quizzes[i].Quiz

			var existing models.Quiz
			err := tx.NewSelect().Model(&existing).Where("path = ?", q.Path).Scan(ctx)
			switch {
			case errors.Is(err, sql.ErrNoRows):
				if _, err := tx.NewInsert().Model(q).Exec(ctx); err != nil {
					return err
				}
				resolved[q.Path] = q.UUID
			case err != nil:
				return err
			case existing.Checksum == q.Checksum:
				resolved[q.Path] = existing.UUID
			default:
				return carefulness.Conflict{Conflictor: q.Path}
			}
		} // multiple round-trips. The code has fallen

		if _, err := tx.NewInsert().Model(&st.Test).Exec(ctx); err != nil {
			return err
		}

		pairs := make([]models.TestsQuizzes, 0, len(st.QuizPaths))
		for pos, p := range st.QuizPaths {
			id, ok := resolved[p]
			if !ok {
				return fmt.Errorf("test references unknown quiz %q", p)
			}
			pairs = append(pairs, models.TestsQuizzes{
				TestUUID: st.Test.UUID, // read AFTER insert
				QuizUUID: id,
				Position: pos,
			})
		}
		if len(pairs) == 0 {
			return nil
		}
		_, err := tx.NewInsert().Model(&pairs).Exec(ctx)
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
