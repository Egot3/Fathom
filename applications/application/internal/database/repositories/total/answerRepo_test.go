package total_test

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json/v2"
	"fmt"
	"math"
	mrand "math/rand/v2"
	"testing"
	"time"

	"github.com/egot3/fathom/internal/database/repositories"
	"github.com/egot3/fathom/internal/database/repositories/total"
	"github.com/egot3/fathom/internal/models"
	"github.com/egot3/fathom/internal/quiz"
	"github.com/egot3/fathom/internal/testutils"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func NewInjectorWithTestRepo(t testing.TB) do.Injector {
	t.Helper()

	i := testutils.NewTestInjector(t)

	do.Provide(i, total.NewTotalRepository)

	return i
}

func TestAnswer_Set(t *testing.T) {
	i := NewInjectorWithTestRepo(t)
	r := do.MustInvoke[total.TotalRepository](i)
	db := do.MustInvoke[*bun.DB](i)

	testM := models.Test{Name: rand.Text()}
	err := db.NewInsert().Model(&testM).Returning("*").Scan(t.Context())
	require.NoError(t, err)

	var quiz models.Quiz = models.Quiz{Path: fmt.Sprintf("/path/to/%v.md", rand.Text()), Checksum: [8]byte{}, Score: 1}
	err = db.NewInsert().Model(&quiz).Returning("*").Scan(t.Context())
	require.NoError(t, err)

	_, err = db.NewInsert().Model(&models.TestsQuizzes{TestUUID: testM.UUID, QuizUUID: quiz.UUID}).Exec(t.Context())
	require.NoError(t, err)

	name := rand.Text()
	var groupUUID uuid.UUID
	err = db.NewInsert().Model(&models.Group{Name: name}).Returning("uuid").Scan(t.Context(), &groupUUID)
	require.NoError(t, err)

	var userUUID uuid.UUID
	err = db.NewInsert().Model(&models.User{Nickname: rand.Text(), PasswordHash: []byte{}}).
		Returning("uuid").
		Scan(t.Context(), &userUUID)
	require.NoError(t, err)
	require.NotEmpty(t, userUUID)

	_, err = db.NewInsert().Model(&models.GroupsUsers{GroupUUID: groupUUID, UserUUID: userUUID}).Exec(t.Context())

	testCases := []struct {
		desc      string
		testUUID  uuid.UUID
		userUUID  uuid.UUID
		groupUUID uuid.UUID
		quizUUID  uuid.UUID
		expectErr bool
	}{
		{
			desc:      "Valid set",
			testUUID:  testM.UUID,
			quizUUID:  quiz.UUID,
			userUUID:  userUUID,
			groupUUID: groupUUID,
			expectErr: false,
		},
		{
			desc:      "Invalid user",
			testUUID:  testM.UUID,
			quizUUID:  quiz.UUID,
			userUUID:  uuid.Nil,
			groupUUID: groupUUID,
			expectErr: true,
		},
		{
			desc:      "Invalid group",
			testUUID:  testM.UUID,
			quizUUID:  quiz.UUID,
			userUUID:  userUUID,
			groupUUID: uuid.Nil,
			expectErr: true,
		},
		{
			desc:      "Invalid group and user", // testing those pairs as they are only binded between eachother
			testUUID:  testM.UUID,
			quizUUID:  quiz.UUID,
			userUUID:  uuid.Nil,
			groupUUID: uuid.Nil,
			expectErr: true,
		},
		{
			desc:      "Invalid all",
			testUUID:  uuid.Nil,
			quizUUID:  uuid.Nil,
			userUUID:  uuid.Nil,
			groupUUID: uuid.Nil,
			expectErr: true,
		},
	} // no test and quiz as it's managed at upper layer
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Cleanup(func() {
				_, err := db.NewTruncateTable().Model((*models.Answer)(nil)).Exec(context.Background())
				require.NoError(t, err)
			})

			err := r.SetAnswer(t.Context(), tC.testUUID, tC.groupUUID, tC.userUUID, tC.quizUUID, "", 1)
			if tC.expectErr {
				require.Error(t, err, err.Error())
				require.ErrorIs(t, err, sql.ErrNoRows, err.Error())
				return
			}
			require.NoError(t, err)
		})
	}

	t.Run("Upsertion", func(t *testing.T) {
		t.Cleanup(func() {
			_, err := db.NewTruncateTable().Model((*models.Answer)(nil)).Exec(context.Background())
			require.NoError(t, err)
		})

		err := r.SetAnswer(t.Context(), testM.UUID, groupUUID, userUUID, quiz.UUID, "", 1)

		require.NoError(t, err)
		var insertion models.Answer = models.Answer{}
		err = db.NewSelect().Model(&insertion).
			Where("test_uuid = ?", testM.UUID).
			Where("group_uuid = ?", groupUUID).
			Where("user_uuid = ?", userUUID).
			Where("quiz_uuid = ?", quiz.UUID).
			Scan(t.Context())
		require.NoError(t, err)

		err = r.SetAnswer(t.Context(), testM.UUID, groupUUID, userUUID, quiz.UUID, "", 1)

		require.NoError(t, err)

		var updation models.Answer = models.Answer{}
		err = db.NewSelect().Model(&insertion).
			Where("test_uuid = ?", testM.UUID).
			Where("group_uuid = ?", groupUUID).
			Where("user_uuid = ?", userUUID).
			Where("quiz_uuid = ?", quiz.UUID).
			Scan(t.Context())
		require.NoError(t, err)

		require.NotEqual(t, insertion.AnsweredAt, updation.AnsweredAt)
	})
}

func TestAnswer_Get(t *testing.T) {
	i := NewInjectorWithTestRepo(t)
	r := do.MustInvoke[total.TotalRepository](i)
	db := do.MustInvoke[*bun.DB](i)

	testM := models.Test{Name: rand.Text()}
	err := db.NewInsert().Model(&testM).Returning("*").Scan(t.Context())
	require.NoError(t, err)

	var quiz models.Quiz = models.Quiz{Path: fmt.Sprintf("/path/to/%v.md", rand.Text()), Checksum: [8]byte{}, Score: 1}
	err = db.NewInsert().Model(&quiz).Returning("*").Scan(t.Context())
	require.NoError(t, err)

	_, err = db.NewInsert().Model(&models.TestsQuizzes{TestUUID: testM.UUID, QuizUUID: quiz.UUID}).Exec(t.Context())
	require.NoError(t, err)

	name := rand.Text()
	groupUUID := uuid.UUID{}
	err = db.NewInsert().Model(&models.Group{Name: name}).Returning("uuid").Scan(t.Context(), &groupUUID)
	require.NoError(t, err)

	var userUUID uuid.UUID

	err = db.NewInsert().Model(&models.User{Nickname: rand.Text(), PasswordHash: []byte{}}).
		Returning("uuid").
		Scan(t.Context(), &userUUID)
	require.NoError(t, err)
	require.NotEmpty(t, userUUID)

	_, err = db.NewInsert().Model(&models.GroupsUsers{GroupUUID: groupUUID, UserUUID: userUUID}).Exec(t.Context())

	score := mrand.Float32() * 256
	totalValue := rand.Text()
	_, err = db.NewInsert().Model(&models.Answer{
		GroupUUID:   groupUUID,
		TestUUID:    testM.UUID,
		UserUUID:    userUUID,
		QuizUUID:    quiz.UUID,
		AnswerValue: totalValue,
		Score:       score,
	}).Exec(t.Context())
	require.NoError(t, err)

	t.Run("Score suite", func(t *testing.T) {
		t.Run("Of known", func(t *testing.T) {
			scoreR, err := r.AnswerScore(t.Context(), userUUID, testM.UUID, groupUUID, quiz.UUID)
			require.NoError(t, err)
			require.EqualValues(t, score, scoreR)
		})
		t.Run("Of unknown", func(t *testing.T) {
			scoreR, err := r.AnswerScore(t.Context(), uuid.Nil, uuid.Nil, uuid.Nil, uuid.Nil)
			require.Error(t, err)
			require.EqualValues(t, 0, scoreR)
		})
	})
	t.Run("Value suite", func(t *testing.T) {
		t.Run("Of known", func(t *testing.T) {
			totalR, err := r.Answer(t.Context(), userUUID, testM.UUID, groupUUID, quiz.UUID)
			require.NoError(t, err)
			require.Equal(t, totalValue, totalR)
		})
		t.Run("Of unknown", func(t *testing.T) {
			totalR, err := r.Answer(t.Context(), uuid.Nil, uuid.Nil, uuid.Nil, uuid.Nil)
			require.Error(t, err)
			require.Equal(t, "", totalR)
		})
	})
}

func TestAnswer_Totalization(t *testing.T) {
	i := NewInjectorWithTestRepo(t)
	r := do.MustInvoke[total.TotalRepository](i)
	db := do.MustInvoke[*bun.DB](i)

	testM := models.Test{Name: rand.Text()}
	err := db.NewInsert().Model(&testM).Returning("*").Scan(t.Context())
	require.NoError(t, err)

	var quiz models.Quiz = models.Quiz{Path: fmt.Sprintf("/path/to/%v.md", rand.Text()), Checksum: [8]byte{}, Score: 1}
	err = db.NewInsert().Model(&quiz).Returning("*").Scan(t.Context())
	require.NoError(t, err)

	_, err = db.NewInsert().Model(&models.TestsQuizzes{TestUUID: testM.UUID, QuizUUID: quiz.UUID}).Exec(t.Context())
	require.NoError(t, err)

	name := rand.Text()
	groupUUID := uuid.UUID{}
	err = db.NewInsert().Model(&models.Group{Name: name}).Returning("uuid").Scan(t.Context(), &groupUUID)
	require.NoError(t, err)

	var userUUID uuid.UUID

	err = db.NewInsert().Model(&models.User{Nickname: rand.Text(), PasswordHash: []byte{}}).
		Returning("uuid").
		Scan(t.Context(), &userUUID)
	require.NoError(t, err)
	require.NotEmpty(t, userUUID)

	_, err = db.NewInsert().Model(&models.GroupsUsers{GroupUUID: groupUUID, UserUUID: userUUID}).Exec(t.Context())

	score := mrand.Float32() * 256
	_, err = db.NewInsert().Model(&models.Answer{
		GroupUUID:   groupUUID,
		TestUUID:    testM.UUID,
		UserUUID:    userUUID,
		QuizUUID:    quiz.UUID,
		AnswerValue: rand.Text(),
		Score:       score,
	}).Exec(t.Context())
	require.NoError(t, err)

	t.Run("Totalize known", func(t *testing.T) {
		t.Cleanup(func() {
			_, err := db.NewTruncateTable().Model((*models.UserGroupsTests)(nil)).Exec(context.Background())
			require.NoError(t, err)
		})

		err := r.Totalize(t.Context(), userUUID, testM.UUID, groupUUID)
		require.NoError(t, err)

		var scoreR float64
		err = db.NewSelect().Model(&models.UserGroupsTests{
			UserUUID:  userUUID,
			GroupUUID: groupUUID,
			TestUUID:  testM.UUID,
		}).Column("score").Scan(t.Context(), &scoreR)

		require.NoError(t, err)
		require.Condition(t, func() (success bool) {
			return math.Abs(scoreR-float64(score)) < 0.1
		})
	})

	t.Run("Totalize unknown", func(t *testing.T) {
		t.Cleanup(func() {
			_, err := db.NewTruncateTable().Model((*models.UserGroupsTests)(nil)).Exec(context.Background())
			require.NoError(t, err)
		})

		err := r.Totalize(t.Context(), uuid.Nil, testM.UUID, groupUUID)
		require.Error(t, err)

		c, err := db.NewSelect().Model((*models.UserGroupsTests)(nil)).
			Column("score").Count(t.Context())

		require.NoError(t, err)
		require.Equal(t, 0, c)
	})

	t.Run("Re-totalize known", func(t *testing.T) {
		t.Cleanup(func() {
			_, err := db.NewTruncateTable().Model((*models.UserGroupsTests)(nil)).Exec(context.Background())
			require.NoError(t, err)
		})

		err := r.Totalize(t.Context(), userUUID, testM.UUID, groupUUID)
		require.NoError(t, err)

		var scoreR float64
		err = db.NewSelect().Model(&models.UserGroupsTests{
			UserUUID:  userUUID,
			GroupUUID: groupUUID,
			TestUUID:  testM.UUID,
		}).Column("score").Scan(t.Context(), &scoreR)

		require.NoError(t, err)
		require.Condition(t, func() (success bool) {
			return math.Abs(scoreR-float64(score)) < 0.1
		})

		scoreN := mrand.Float32() * 256
		var quiz models.Quiz = models.Quiz{Path: fmt.Sprintf("/path/to/%v.md", rand.Text()), Checksum: [8]byte{}, Score: 1}
		err = db.NewInsert().Model(&quiz).Returning("*").Scan(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.TestsQuizzes{TestUUID: testM.UUID, QuizUUID: quiz.UUID}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.Answer{
			GroupUUID:   groupUUID,
			TestUUID:    testM.UUID,
			UserUUID:    userUUID,
			QuizUUID:    quiz.UUID,
			AnswerValue: rand.Text(),
			Score:       scoreN,
		}).Exec(t.Context())
		require.NoError(t, err)

		err = r.Totalize(t.Context(), userUUID, testM.UUID, groupUUID)
		require.NoError(t, err)

		var scoreR2 float32
		err = db.NewSelect().Model(&models.UserGroupsTests{
			UserUUID:  userUUID,
			GroupUUID: groupUUID,
			TestUUID:  testM.UUID,
		}).Column("score").Scan(t.Context(), &scoreR2)

		require.NoError(t, err)
		require.InDelta(t, score+scoreN, scoreR2, 1e-4)
	})
}

func TestAnswer_Totals(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		i := testutils.NewTestInjector(t, repositories.RepositoryPackage)
		r := do.MustInvoke[total.TotalRepository](i)
		db := do.MustInvoke[*bun.DB](i)

		testM := models.Test{Name: rand.Text()}
		err := db.NewInsert().Model(&testM).Returning("*").Scan(t.Context())
		require.NoError(t, err)

		var quiz models.Quiz = models.Quiz{Path: fmt.Sprintf("/path/to/%v.md", rand.Text()), Checksum: [8]byte{}, Score: mrand.IntN(32)}
		err = db.NewInsert().Model(&quiz).Returning("*").Scan(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.TestsQuizzes{TestUUID: testM.UUID, QuizUUID: quiz.UUID}).Exec(t.Context())
		require.NoError(t, err)

		name := rand.Text()
		groupUUID := uuid.UUID{}
		err = db.NewInsert().Model(&models.Group{Name: name}).Returning("uuid").Scan(t.Context(), &groupUUID)
		require.NoError(t, err)

		var userUUID uuid.UUID

		err = db.NewInsert().Model(&models.User{Nickname: rand.Text(), PasswordHash: []byte{}}).
			Returning("uuid").
			Scan(t.Context(), &userUUID)
		require.NoError(t, err)
		require.NotEmpty(t, userUUID)

		_, err = db.NewInsert().Model(&models.GroupsUsers{GroupUUID: groupUUID, UserUUID: userUUID}).Exec(t.Context())

		score := mrand.Float64() * 256
		_, err = db.NewInsert().Model(&models.UserGroupsTests{
			TestUUID:  testM.UUID,
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
			Score:     score,
		}).Exec(t.Context())
		require.NoError(t, err)

		totals, totalN, err := r.ListTotals(t.Context(), 0, 5, userUUID, testM.UUID, groupUUID)
		require.NoError(t, err)
		require.EqualValues(t, 1, totalN)
		require.Len(t, totals, 1)

		t.Logf("%+v", totals)
	})
	t.Run("None", func(t *testing.T) {
		i := testutils.NewTestInjector(t, repositories.RepositoryPackage)
		r := do.MustInvoke[total.TotalRepository](i)

		to, _, err := r.ListTotals(t.Context(), 0, 5, uuid.Nil, uuid.Nil, uuid.Nil)
		require.ErrorIs(t, err, sql.ErrNoRows, fmt.Sprintf("%+v", to))
	})
	t.Run("Multiple", func(t *testing.T) {
		i := testutils.NewTestInjector(t, repositories.RepositoryPackage)
		r := do.MustInvoke[total.TotalRepository](i)
		db := do.MustInvoke[*bun.DB](i)

		amount := mrand.IntN(32) + 1
		ugts := make([]models.UserGroupsTests, 0, amount)

		for range amount {
			testM := models.Test{Name: rand.Text()}
			err := db.NewInsert().Model(&testM).Returning("*").Scan(t.Context())
			require.NoError(t, err)

			var quiz models.Quiz = models.Quiz{Path: fmt.Sprintf("/path/to/%v.md", rand.Text()), Checksum: [8]byte{}, Score: mrand.Int() + 1}
			err = db.NewInsert().Model(&quiz).Returning("*").Scan(t.Context())
			require.NoError(t, err)

			_, err = db.NewInsert().Model(&models.TestsQuizzes{TestUUID: testM.UUID, QuizUUID: quiz.UUID}).Exec(t.Context())
			require.NoError(t, err)

			name := rand.Text()
			groupUUID := uuid.UUID{}
			err = db.NewInsert().Model(&models.Group{Name: name}).Returning("uuid").Scan(t.Context(), &groupUUID)
			require.NoError(t, err)

			var userUUID uuid.UUID

			err = db.NewInsert().Model(&models.User{Nickname: rand.Text(), PasswordHash: []byte{}}).
				Returning("uuid").
				Scan(t.Context(), &userUUID)
			require.NoError(t, err)
			require.NotEmpty(t, userUUID)

			_, err = db.NewInsert().Model(&models.GroupsUsers{GroupUUID: groupUUID, UserUUID: userUUID}).Exec(t.Context())

			score := mrand.Float64() * 256
			ugts = append(ugts, models.UserGroupsTests{
				UserUUID:  userUUID,
				GroupUUID: groupUUID,
				TestUUID:  testM.UUID,
				Score:     score,
			})
		}

		_, err := db.NewInsert().Model(&ugts).Exec(t.Context())
		require.NoError(t, err)

		totals, N, err := r.ListTotals(t.Context(), 0, int(amount+1), uuid.Nil, uuid.Nil, uuid.Nil)
		require.NoError(t, err)
		require.EqualValues(t, N, amount)
		require.Len(t, totals, int(amount))
	})
}

func TestAnswer_AnswersInTest(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		t.Run("Single", func(t *testing.T) {
			i := NewInjectorWithTestRepo(t)
			r := do.MustInvoke[total.TotalRepository](i)
			db := do.MustInvoke[*bun.DB](i)

			testM := models.Test{Name: rand.Text()}
			err := db.NewInsert().Model(&testM).Returning("*").Scan(t.Context())
			require.NoError(t, err)

			qans := quiz.QuizAnswers{Input: &quiz.AnswerInput{Input: rand.Text()}}
			qscore := mrand.IntN(256) + 1
			quizAns, err := json.Marshal(qans)
			require.NoError(t, err)
			q := models.Quiz{Path: fmt.Sprintf("/path/to/%v.md", rand.Text()), Checksum: [8]byte{}, Score: int(qscore), CorrectAnswer: string(quizAns)}
			err = db.NewInsert().Model(&q).Returning("*").Scan(t.Context())
			require.NoError(t, err)

			_, err = db.NewInsert().Model(&models.TestsQuizzes{TestUUID: testM.UUID, QuizUUID: q.UUID}).Exec(t.Context())
			require.NoError(t, err)

			name := rand.Text()
			groupUUID := uuid.UUID{}
			err = db.NewInsert().Model(&models.Group{Name: name}).Returning("uuid").Scan(t.Context(), &groupUUID)
			require.NoError(t, err)

			var userUUID uuid.UUID

			err = db.NewInsert().Model(&models.User{Nickname: rand.Text(), PasswordHash: []byte{}}).
				Returning("uuid").
				Scan(t.Context(), &userUUID)
			require.NoError(t, err)
			require.NotEmpty(t, userUUID)

			_, err = db.NewInsert().Model(&models.GroupsUsers{GroupUUID: groupUUID, UserUUID: userUUID}).Exec(t.Context())

			ans := quiz.QuizAnswers{Input: &quiz.AnswerInput{Input: rand.Text()}}
			score := mrand.Float32() * 256
			totalValue, err := json.Marshal(ans)
			require.NoError(t, err)
			_, err = db.NewInsert().Model(&models.Answer{
				GroupUUID:   groupUUID,
				TestUUID:    testM.UUID,
				UserUUID:    userUUID,
				QuizUUID:    q.UUID,
				AnswerValue: string(totalValue),
				Score:       score,
			}).Exec(t.Context())
			require.NoError(t, err)

			answers, total, err := r.AnswersInTest(t.Context(), userUUID, testM.UUID, groupUUID, 0, 5)
			require.NoError(t, err)
			require.Equal(t, 1, total)
			require.Len(t, answers, 1)

			require.Equal(t, ans, answers[0].Chosen)
		})
		t.Run("Multiple", func(t *testing.T) {
			i := NewInjectorWithTestRepo(t)
			r := do.MustInvoke[total.TotalRepository](i)
			db := do.MustInvoke[*bun.DB](i)

			testM := models.Test{Name: rand.Text()}
			err := db.NewInsert().Model(&testM).Returning("*").Scan(t.Context())
			require.NoError(t, err)

			name := rand.Text()
			groupUUID := uuid.UUID{}
			err = db.NewInsert().Model(&models.Group{Name: name}).Returning("uuid").Scan(t.Context(), &groupUUID)
			require.NoError(t, err)

			var userUUID uuid.UUID

			err = db.NewInsert().Model(&models.User{Nickname: rand.Text(), PasswordHash: []byte{}}).
				Returning("uuid").
				Scan(t.Context(), &userUUID)
			require.NoError(t, err)
			require.NotEmpty(t, userUUID)

			_, err = db.NewInsert().Model(&models.GroupsUsers{GroupUUID: groupUUID, UserUUID: userUUID}).Exec(t.Context())
			require.NoError(t, err)

			amount := mrand.Int64N(256) + 10 // when 2 minutes left, brain works faster by coef 5/t, where t is time in minutes
			var latest quiz.QuizAnswers
			for i := range amount {
				qans := quiz.QuizAnswers{Input: &quiz.AnswerInput{Input: rand.Text()}}
				qscore := mrand.IntN(256) + 1
				quizAns, err := json.Marshal(qans)
				require.NoError(t, err)
				q := models.Quiz{
					Path:          fmt.Sprintf("/path/to/%v.md", rand.Text()),
					Checksum:      [8]byte{},
					Score:         (qscore),
					CorrectAnswer: string(quizAns),
				}
				err = db.NewInsert().Model(&q).Returning("*").Scan(t.Context())
				require.NoError(t, err)

				_, err = db.NewInsert().Model(&models.TestsQuizzes{TestUUID: testM.UUID, QuizUUID: q.UUID}).Exec(t.Context())
				require.NoError(t, err)

				ans := quiz.QuizAnswers{Input: &quiz.AnswerInput{Input: rand.Text()}}
				score := mrand.Float32() * 256
				totalValue, err := json.Marshal(ans)
				require.NoError(t, err)
				_, err = db.NewInsert().Model(&models.Answer{
					GroupUUID:   groupUUID,
					TestUUID:    testM.UUID,
					UserUUID:    userUUID,
					QuizUUID:    q.UUID,
					AnswerValue: string(totalValue),
					Score:       score,
					AnsweredAt:  time.Now().Add(time.Duration(i) * time.Second),
				}).Exec(t.Context())
				require.NoError(t, err)

				latest = ans
			}

			answers, total, err := r.AnswersInTest(t.Context(), userUUID, testM.UUID, groupUUID, 0, int(amount))
			require.NoError(t, err)
			require.EqualValues(t, amount, total)
			require.Len(t, answers, int(amount))

			require.Equal(t, latest, answers[0].Chosen)
		})
	})
}

func TestAnswer_LatestAnswers(t *testing.T) {
	i := NewInjectorWithTestRepo(t)
	t.Run("Valid", func(t *testing.T) {
		t.Run("Single", func(t *testing.T) {
			i := i.Scope(t.Name())
			r := do.MustInvoke[total.TotalRepository](i)
			db := do.MustInvoke[*bun.DB](i)

			testM := models.Test{Name: rand.Text()}
			err := db.NewInsert().Model(&testM).Returning("*").Scan(t.Context())
			require.NoError(t, err)

			var quiz models.Quiz = models.Quiz{Path: fmt.Sprintf("/path/to/%v.md", rand.Text()), Checksum: [8]byte{}, Score: 1}
			err = db.NewInsert().Model(&quiz).Returning("*").Scan(t.Context())
			require.NoError(t, err)

			_, err = db.NewInsert().Model(&models.TestsQuizzes{TestUUID: testM.UUID, QuizUUID: quiz.UUID}).Exec(t.Context())
			require.NoError(t, err)

			name := rand.Text()
			groupUUID := uuid.UUID{}
			err = db.NewInsert().Model(&models.Group{Name: name}).Returning("uuid").Scan(t.Context(), &groupUUID)
			require.NoError(t, err)

			var userUUID uuid.UUID
			err = db.NewInsert().Model(&models.User{Nickname: rand.Text(), PasswordHash: []byte{}}).
				Returning("uuid").
				Scan(t.Context(), &userUUID)
			require.NoError(t, err)
			require.NotEmpty(t, userUUID)

			_, err = db.NewInsert().Model(&models.GroupsUsers{GroupUUID: groupUUID, UserUUID: userUUID}).Exec(t.Context())

			score := mrand.Float32() * 256
			totalValue := rand.Text()
			_, err = db.NewInsert().Model(&models.Answer{
				GroupUUID:   groupUUID,
				TestUUID:    testM.UUID,
				UserUUID:    userUUID,
				QuizUUID:    quiz.UUID,
				AnswerValue: totalValue,
				Score:       score,
			}).Exec(t.Context())
			require.NoError(t, err)

			m, err := r.LatestAnswers(t.Context(), userUUID, testM.UUID, groupUUID, time.Now().Add(-1*time.Second), uuid.UUIDs{quiz.UUID})
			require.NoError(t, err)

			require.Equal(t, totalValue, m[quiz.UUID])
		})
		t.Run(t.Name(), func(t *testing.T) {
			i := i.Scope(t.Name())
			r := do.MustInvoke[total.TotalRepository](i)
			db := do.MustInvoke[*bun.DB](i)

			testM := models.Test{Name: rand.Text()}
			err := db.NewInsert().Model(&testM).Returning("*").Scan(t.Context())
			require.NoError(t, err)

			var quiz models.Quiz = models.Quiz{Path: fmt.Sprintf("/path/to/%v.md", rand.Text()), Checksum: [8]byte{}, Score: 1}
			err = db.NewInsert().Model(&quiz).Returning("*").Scan(t.Context())
			require.NoError(t, err)

			_, err = db.NewInsert().Model(&models.TestsQuizzes{TestUUID: testM.UUID, QuizUUID: quiz.UUID}).Exec(t.Context())
			require.NoError(t, err)

			name := rand.Text()
			groupUUID := uuid.UUID{}
			err = db.NewInsert().Model(&models.Group{Name: name}).Returning("uuid").Scan(t.Context(), &groupUUID)
			require.NoError(t, err)

			var userUUID uuid.UUID
			err = db.NewInsert().Model(&models.User{Nickname: rand.Text(), PasswordHash: []byte{}}).
				Returning("uuid").
				Scan(t.Context(), &userUUID)
			require.NoError(t, err)
			require.NotEmpty(t, userUUID)

			_, err = db.NewInsert().Model(&models.GroupsUsers{GroupUUID: groupUUID, UserUUID: userUUID}).Exec(t.Context())

			var tValue string
			for range mrand.IntN(32) + 1 {
				score := mrand.Float32() * 256
				totalValue := rand.Text()
				_, err = db.NewInsert().Model(&models.Answer{
					GroupUUID:   groupUUID,
					TestUUID:    testM.UUID,
					UserUUID:    userUUID,
					QuizUUID:    quiz.UUID,
					AnswerValue: totalValue,
					Score:       score,
				}).Exec(t.Context())
				require.NoError(t, err)
				tValue = totalValue
			}

			m, err := r.LatestAnswers(t.Context(), userUUID, testM.UUID, groupUUID, time.Now().Add(-1*time.Second), uuid.UUIDs{quiz.UUID})
			require.NoError(t, err)

			require.Equal(t, tValue, m[quiz.UUID])
		})
	})
}
