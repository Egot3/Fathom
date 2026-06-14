package test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

var ErrQuizNotInTest = errors.New("quiz is not in the runner")

type NotInTestError struct {
	Pathes []string
}

func (e *NotInTestError) Error() string {
	return fmt.Sprintf("%d quizzes are not in test", len(e.Pathes))
}

func (e *NotInTestError) Is(target error) bool {
	return target == ErrQuizNotInTest
}

type bunTestRepository struct {
	db *bun.DB
}

func NewTestRepository(i do.Injector) (TestRepository, error) {
	db := do.MustInvoke[*bun.DB](i)
	return &bunTestRepository{db: db}, nil
}

func (r *bunTestRepository) CreateTest(ctx context.Context, name string) error {
	_, err := r.db.NewInsert().Model(&models.Test{Name: name}).Exec(ctx)
	return err
}

func (r *bunTestRepository) BundleQuizzesToTest(ctx context.Context, testUUID uuid.UUID, pathes []string) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, path := range pathes {
			_, err := tx.NewInsert().Model(&models.TestsQuizzies{TestUUID: testUUID, QuizPath: path}).Exec(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *bunTestRepository) PruneQuizzesFromTest(ctx context.Context, testUUID uuid.UUID, pathes []string) error {
	notFound := make([]string, 0)
	err := r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, path := range pathes {
			_, err := tx.NewDelete().Model(&models.TestsQuizzies{QuizPath: path, TestUUID: testUUID}).WherePK().Exec(ctx)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					notFound = append(notFound, path)
				}
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	if len(notFound) > 0 {
		return &NotInTestError{Pathes: notFound}
	}

	return nil
}

func (r *bunTestRepository) Test(ctx context.Context, UUID uuid.UUID) (*models.Test, error) {
	var test = &models.Test{UUID: UUID}
	err := r.db.NewSelect().Model(&test).WherePK().Relation("Quiz").Scan(ctx)
	if err != nil {
		return nil, err
	}

	return test, nil
}

func (r *bunTestRepository) DeleteTest(ctx context.Context, UUID uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&models.Test{UUID: UUID}).WherePK().Exec(ctx)
	return err
}

func (r *bunTestRepository) UpdateTest(ctx context.Context, UUID uuid.UUID, name string) error {
	_, err := r.db.NewUpdate().Model(&models.Test{UUID: UUID, Name: name}).WherePK().Exec(ctx)
	return err
}
