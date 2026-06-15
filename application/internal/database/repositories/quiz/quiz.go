package quiz

import (
	"context"
	"database/sql"

	"github.com/egot3/fathom/internal/models"
	"github.com/samber/do/v2"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
)

type bunQuizRepository struct {
	db *bun.DB
}

func NewQuizRepository(i do.Injector) (QuizRepository, error) {
	db := do.MustInvoke[*bun.DB](i)
	return &bunQuizRepository{db: db}, nil
}

func (r *bunQuizRepository) RegisterQuiz(ctx context.Context, path string, checksum []byte, score int) error {
	_, err := r.db.NewInsert().Model(&models.Quiz{Path: path, Checksum: checksum, Score: score}).Exec(ctx)
	return err
}

func (r *bunQuizRepository) DeallocateQuiz(ctx context.Context, path string) error {
	res, err := r.db.NewDelete().Model(&models.Quiz{Path: path}).WherePK().Exec(ctx)
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

func (r *bunQuizRepository) ListQuizzes(ctx context.Context, page, size int) ([]string, int, error) {
	var quizzes []models.Quiz
	total, err := r.db.NewSelect().Model(&quizzes).Offset(page*size).Limit(page).OrderBy("path", bun.OrderDesc).ScanAndCount(ctx)
	if err != nil {
		return nil, total, err
	}

	return lo.Map(quizzes, func(quiz models.Quiz, _ int) string { return quiz.Path }), total, nil
}

func (r *bunQuizRepository) CheckRegistered(ctx context.Context, path string) (bool, error) {
	return r.db.NewSelect().Model(&models.Quiz{Path: path}).WherePK().Exists(ctx)
}

func (r *bunQuizRepository) CheckIntegrity(ctx context.Context, path string, checksum []byte) (bool, error) {
	return r.db.NewSelect().Model(&models.Quiz{Path: path}).WherePK().Where("checksum = ?", checksum).Exists(ctx)
}
