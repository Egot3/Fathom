package test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/egot3/fathom/internal/carefulness"
	exportutlis "github.com/egot3/fathom/internal/exportUtlis"
	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
)

var ErrQuizNotInTest = errors.New("quiz is not in the test")

type NotInTestError struct {
	Count int
}

func (e *NotInTestError) Error() string {
	return fmt.Sprintf("%d quizzes are not in test", e.Count)
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

func (r *bunTestRepository) TestPathes(ctx context.Context, UUIDs uuid.UUIDs) ([]string, error) {
	pathes := make([]string, len(UUIDs))
	err := r.db.NewSelect().Model((*models.Quiz)(nil)).
		Column("path").
		Where("uuid IN (?)", bun.List(UUIDs)).
		Scan(ctx, pathes)
	if err != nil {
		return nil, err
	}

	if len(UUIDs) != len(pathes) {
		return nil, sql.ErrNoRows
	}

	return pathes, nil
}

func (r *bunTestRepository) CreateTest(ctx context.Context, name string) (*models.Test, error) {
	var test = models.Test{Name: name}
	err := r.db.NewInsert().Model(&test).Returning("*").Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &carefulness.Conflict{Conflictor: "name"}
		}
		return nil, err
	}

	return &test, nil
}

func (r *bunTestRepository) BundleQuizzesToTest(ctx context.Context, testUUID uuid.UUID, quizUUIDs uuid.UUIDs) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		var testQuizzes []models.TestsQuizzies = lo.Map(quizUUIDs, func(u uuid.UUID, pos int) models.TestsQuizzies {
			return models.TestsQuizzies{TestUUID: testUUID, Position: pos, QuizUUID: u}
		})
		_, err := tx.NewInsert().Model(&testQuizzes).Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}

/* func (r *bunTestRepository) BundleQuizzesToTest(ctx context.Context, testUUID uuid.UUID, pathes []string) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		var testQuizzes []models.TestsQuizzies = make([]models.TestsQuizzies, len(pathes))
		for pos, path := range pathes {
			testQuizzes[pos] = models.TestsQuizzies{Position: pos, TestUUID: testUUID, QuizPath: path}
		}
		_, err := tx.NewInsert().Model(&testQuizzes).Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	})
} */

// old
/* func (r *bunTestRepository) PruneQuizzesFromTest(ctx context.Context, testUUID uuid.UUID, pathes []string) error {
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
		return &NotInTestError{Count: len(notFound)}
	}

	return nil
} */

// new(alpha)
func (r *bunTestRepository) PruneQuizzesFromTest(ctx context.Context, testUUID uuid.UUID, quizUUIDs uuid.UUIDs) error {
	notFound := 0

	res, err := r.db.NewDelete().Model((*models.TestsQuizzies)(nil)).
		Where("test_uuid = ?", testUUID).Where("quiz_uuid IN (?)", bun.List(quizUUIDs)).
		Exec(ctx)
	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c == 0 {
		return sql.ErrNoRows
	}

	notFound = len(quizUUIDs) - int(c)

	if notFound > 0 {
		return &NotInTestError{Count: notFound}
	}

	return nil
}

func (r *bunTestRepository) Test(ctx context.Context, UUID uuid.UUID) (*models.Test, error) {
	var test = models.Test{UUID: UUID}
	err := r.db.NewSelect().Model(&test).WherePK().Relation("Quizzes").Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &test, nil
}

func (r *bunTestRepository) DeleteTest(ctx context.Context, UUID uuid.UUID) error {
	res, err := r.db.NewDelete().Model(&models.Test{UUID: UUID}).WherePK().Exec(ctx)

	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *bunTestRepository) UpdateTest(ctx context.Context, UUID uuid.UUID, name string) error {
	e, err := r.db.NewSelect().Model((*models.Test)(nil)).Where("name = ?", name).Exists(ctx)
	if err != nil {
		return err
	}
	if e {
		return carefulness.Conflict{Conflictor: "name"}
	}

	res, err := r.db.NewUpdate().Model(&models.Test{UUID: UUID, Name: name}).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *bunTestRepository) ListTests(ctx context.Context, page, size int) ([]models.Test, int, error) {
	tests := make([]models.Test, size)
	total, err := r.db.NewSelect().Model(&tests).OrderBy("uuid", bun.OrderAsc).
		Limit(size).Offset(page * size).ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	return tests, total, nil
}

func (r *bunTestRepository) ExistsByUUID(ctx context.Context, testUUID uuid.UUID) (bool, error) {
	return r.db.NewSelect().Model((*models.Test)(nil)).
		Where("uuid = ?", testUUID).Exists(ctx)
}

func (r *bunTestRepository) ImportTest(ctx context.Context, test exportutlis.YamlTest) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(&models.Test{
			UUID: test.UUID,
			Name: test.Name,
		}).Exec(ctx)
		if err != nil {
			return err
		}

		for _, q := range test.Quizzes {
			_, err := tx.NewInsert().Model(&models.TestsQuizzies{
				QuizUUID: q.UUID,
				TestUUID: test.UUID,
			}).Exec(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
