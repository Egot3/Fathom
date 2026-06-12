package quiz

import (
	"context"

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

func (r *bunQuizRepository) RegisterQuiz(ctx context.Context, path string) error {
	_, err := r.db.NewInsert().Model(&models.Quiz{Path: path}).Exec(ctx)
	return err
}

func (r *bunQuizRepository) DeallocateQuiz(ctx context.Context, path string) error {
	_, err := r.db.NewDelete().Model(&models.Quiz{Path: path}).WherePK().Exec(ctx)
	return err
}

func (r *bunQuizRepository) ListQuizzes(ctx context.Context, page, size int) ([]string, int, error) {
	var quizzes []models.Quiz
	total, err := r.db.NewSelect().Model(&quizzes).Offset(page*size).Limit(page).OrderBy("path", bun.OrderDesc).ScanAndCount(ctx)
	if err != nil {
		return nil, total, err
	}

	return lo.Map(quizzes, func(quiz models.Quiz, _ int) string { return quiz.Path }), total, nil
}
