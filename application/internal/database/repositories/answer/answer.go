package answer

import (
	"context"
	"database/sql"

	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunAnswerRepository struct {
	db *bun.DB
}

func NewAnswerRepository(i do.Injector) (AnswerRepository, error) {
	db := do.MustInvoke[*bun.DB](i)
	return &bunAnswerRepository{db: db}, nil
} //created all of it just to live a good life in tests

func (r *bunAnswerRepository) SetAnswer(ctx context.Context, answer models.Answer) error {
	_, err := r.db.NewInsert().Model(&answer).Exec(ctx)
	return err
}

func (r *bunAnswerRepository) Totalize(ctx context.Context, userUUID, testUUID, groupUUID uuid.UUID) error {
	var userTotal uint64
	var quizPathes []string
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		err := tx.NewSelect().TableExpr("tests_quizzes AS tq").
			Column("tq.quiz_path").
			Where("tq.test_uuid = ?", testUUID).
			Scan(ctx, &quizPathes)
		if err != nil {
			return err
		}

		if len(quizPathes) == 0 {
			userTotal = 0
		} else {
			latestPerQuiz := tx.NewSelect().
				TableExpr("users_groups_tests_quiz_answers AS inner_a").
				ColumnExpr("MAX(inner_a.answered_at)").
				Where("inner_a.test_uuid = a.test_uuid").
				Where("inner_a.group_uuid = a.group_uuid").
				Where("inner_a.user_uuid = a.user_uuid").
				Where("inner_a.quiz_path = a.quiz_path")

			err = tx.NewSelect().
				TableExpr("users_groups_tests_quiz_answers AS a").
				ColumnExpr("COALESCE(SUM(a.score), 0)").
				Where("a.test_uuid = ?", testUUID).
				Where("a.group_uuid = ?", groupUUID).
				Where("a.user_uuid = ?", userUUID).
				Where("a.quiz_path IN (?)", bun.List(quizPathes)).
				Where("a.answered_at = (?)", latestPerQuiz).
				Scan(ctx, &userTotal)
			if err != nil {
				return err
			}
		}

		_, err = tx.NewInsert().On("CONFLICT (test_uuid, group_uuid, user_uuid) DO UPDATE").
			Set("score = EXCLUDED.score").
			Model(&models.UserGroupsTests{UserUUID: userUUID, GroupUUID: groupUUID, TestUUID: testUUID, Score: userTotal}).
			Exec(ctx)
		return err
	})
}

func (r *bunAnswerRepository) AnswerScore(ctx context.Context, userUUID, testUUID, groupUUID uuid.UUID, quizPath string) (int, error) {
	var score int
	err := r.db.NewSelect().
		Model((*models.Answer)(nil)).
		Where("test_uuid = ?", testUUID).
		Where("group_uuid = ?", groupUUID).
		Where("user_uuid = ?", userUUID).
		Where("quiz_path = ?", quizPath).
		OrderBy("answered_at", bun.OrderDesc).
		Limit(1).
		Column("score").Scan(ctx, &score)
	if err != nil {
		return 0, err
	}

	return score, nil
}

func (r *bunAnswerRepository) Total(ctx context.Context, userUUID, testUUID, groupUUID uuid.UUID) ([]models.UserGroupsTests, error) {
	var totals []models.UserGroupsTests
	err := r.db.NewSelect().Model(&totals).
		Where("test_uuid = ?", testUUID).
		Where("user_uuid = ?", userUUID).
		Where("group_uuid = ?", groupUUID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return totals, nil
}

func (r *bunAnswerRepository) AllTotals(ctx context.Context, userUUID uuid.UUID) ([]models.UserGroupsTests, error) {
	var totals []models.UserGroupsTests
	err := r.db.NewSelect().Model(&totals).
		Where("user_uuid = ?", userUUID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return totals, nil
}

func (r *bunAnswerRepository) TestTotals(ctx context.Context, testUUID uuid.UUID) ([]models.UserGroupsTests, error) {
	var totals []models.UserGroupsTests
	err := r.db.NewSelect().Model(&totals).
		Where("test_uuid = ?", testUUID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return totals, nil
}

func (r *bunAnswerRepository) GroupTestTotals(ctx context.Context, testUUID, groupUUID uuid.UUID) ([]models.UserGroupsTests, error) {
	var totals []models.UserGroupsTests
	err := r.db.NewSelect().Model(&totals).
		Where("test_uuid = ?", testUUID).
		Where("group_uuid = ?", groupUUID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return totals, nil
}
