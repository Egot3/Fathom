package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/importutils"
	"github.com/egot3/fathom/internal/logging"
	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UUIDPair struct {
	ExistingUUID uuid.UUID `bun:"existing_uuid"`
	AddedUUID    uuid.UUID `bun:"added_uuid"`
}

type resolvedUUID struct {
	UUID   uuid.UUID
	Exists bool
}

func (c *chiService) commitTestImport(
	ctx context.Context,
	quizzes []importutils.StagedQuiz,
	st importutils.StagedTest,
) error {
	var placed []string
	err := c.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		resolved := make(map[string]resolvedUUID, len(quizzes))

		for i := range quizzes {
			q := &quizzes[i].Quiz
			relPath := quizzes[i].RelPath

			var existing models.Quiz
			err := tx.NewSelect().Model(&existing).Where("path = ?", q.Path).Scan(ctx)
			switch {
			case errors.Is(err, sql.ErrNoRows):
				if _, err := tx.NewInsert().Model(q).Exec(ctx); err != nil {
					return err
				}
				resolved[relPath] = resolvedUUID{
					UUID:   q.UUID,
					Exists: false,
				}
			case err != nil:
				return err
			case existing.Checksum == q.Checksum:
				resolved[relPath] = resolvedUUID{
					UUID:   existing.UUID,
					Exists: true,
				}
			default:
				return carefulness.Conflict{Conflictor: relPath}
			}
		} // multiple round-trips. The code has fallen

		if _, err := tx.NewInsert().Model(&st.Test).Exec(ctx); err != nil {
			return err
		}

		pairs := make([]models.TestsQuizzes, 0, len(st.QuizPaths))
		for pos, p := range st.QuizPaths {
			pair, ok := resolved[p]
			if !ok {
				return fmt.Errorf("test references unknown quiz %q", p)
			}
			pairs = append(pairs, models.TestsQuizzes{
				TestUUID: st.Test.UUID, // read AFTER insert
				QuizUUID: pair.UUID,
				Position: pos,
			})
		}
		if len(pairs) == 0 {
			return nil
		}
		_, err := tx.NewInsert().Model(&pairs).Exec(ctx)
		if err != nil {
			return err
		}

		for _, s := range quizzes {
			if resolved[s.RelPath].Exists {
				continue
			}
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

func (c *chiService) commitImport(ctx context.Context, staged []importutils.StagedQuiz) error {
	var placed []string

	err := c.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		for i := range staged {
			q := &staged[i].Quiz
			relPath := staged[i].RelPath

			var existing models.Quiz
			err := tx.NewSelect().Model(&existing).Where("path = ?", q.Path).Scan(ctx)
			switch {
			case errors.Is(err, sql.ErrNoRows):
				if _, err := tx.NewInsert().Model(q).Exec(ctx); err != nil {
					return err
				}
			case err != nil:
				return err
			case existing.Checksum == q.Checksum:
				continue
			default:
				return carefulness.Conflict{Conflictor: relPath}
			}
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

func (c *chiService) commitDeletion(ctx context.Context, quizUUID uuid.UUID, forceDelete bool) carefulness.JSONErrorable {
	logger := logging.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
	)
	ctx = logging.WithLogger(ctx, logger)

	err := c.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		var quizPath string
		err := tx.NewSelect().Model(&models.Quiz{UUID: quizUUID}).
			WherePK().Column("path").Scan(ctx, &quizPath)
		if err != nil {
			logger.Error("unable to find quiz's path",
				slog.String("Error", err.Error()),
			)
			return carefulness.JSONError{Err: "couldn't find quiz's path", Status: http.StatusNotFound}
		}

		err = c.quizRepo.DeallocateQuiz(ctx, quizUUID, tx)
		if err != nil {
			logger.Error("unable to deallocate quiz",
				slog.String("Error", err.Error()),
			)
			if errors.Is(err, sql.ErrNoRows) {
				return carefulness.JSONError{Err: "Quiz not found", Status: http.StatusNotFound}
			}
			return carefulness.JSONError{Err: "couldn't deallocate quiz", Status: http.StatusInternalServerError}
		}

		if forceDelete {
			if err != nil {
				logger.Error("unable to retrieve quizPath",
					slog.String("Error", err.Error()),
				)
				return (carefulness.JSONError{Err: "unable to get test's path internally", Status: http.StatusNotFound})
			}

			err = os.Remove(quizPath)
			if err != nil {
				logger.Error("unable to delete quiz by quizPath",
					slog.String("Error", err.Error()),
				)

				return carefulness.JSONError{Err: "unable to delete quiz", Status: http.StatusInternalServerError}
			}
			return nil
		}

		oldRel, err := c.cfg.TurnToRel(quizPath)
		if err != nil {
			logger.Error("couldn't turn old path to rel", slog.String("Error", err.Error()))

			return carefulness.JSONError{Err: "couldn't get old file's rel path", Status: http.StatusInternalServerError}
		}

		deleteDir, err := c.cfg.TurnToAbs("/.deleted")
		if err != nil {
			logger.Error("couldn't turn delete directory to abs", slog.String("Error", err.Error()))

			return carefulness.JSONError{Err: "couldn't ensure delete directory exists", Status: http.StatusInternalServerError}
		}
		newAbs := filepath.Join(deleteDir, oldRel)

		if err := os.MkdirAll(filepath.Dir(newAbs), 0o755); err != nil {
			logger.Error("couldn't ensure delete directory exists", slog.String("Error", err.Error()))

			return carefulness.JSONError{Err: "couldn't ensure delete directory exists", Status: http.StatusInternalServerError}
		}

		err = os.Rename(quizPath, newAbs)
		if err != nil {
			logger.Error("couldn't get rename to abs name", slog.String("Error", err.Error()))

			return carefulness.JSONError{Err: err.Error(), Status: http.StatusInternalServerError}

		}

		return nil
	})

	if err != nil {
		if jerr, ok := errors.AsType[carefulness.JSONErrorable](err); ok {
			return jerr
		}
		panic("not jsonerrorable got from transaction")
	}

	return nil
}
