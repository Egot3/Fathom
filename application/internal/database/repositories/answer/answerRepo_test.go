package answer_test

import (
	"context"
	"crypto/rand"
	"fmt"
	mrand "math/rand/v2"
	"testing"

	"github.com/egot3/fathom/internal/database/repositories/answer"
	"github.com/egot3/fathom/internal/models"
	"github.com/egot3/fathom/internal/testutils"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func NewInjectorWithTestRepo(t testing.TB) do.Injector {
	t.Helper()

	i := testutils.NewTestInjector(t)

	do.Provide(i, answer.NewAnswerRepository)

	return i
}

func RegisterModels(db *bun.DB) {
	db.RegisterModel((*models.TestsQuizzies)(nil))
	db.RegisterModel((*models.GroupsUsers)(nil))
	db.RegisterModel((*models.UserGroupsTests)(nil))
	db.RegisterModel((*models.Answer)(nil))
}

func TestAnswer_Set(t *testing.T) {
	i := NewInjectorWithTestRepo(t)
	r := do.MustInvoke[answer.AnswerRepository](i)
	db := do.MustInvoke[*bun.DB](i)

	RegisterModels(db)

	testM := models.Test{Name: rand.Text()}
	err := db.NewInsert().Model(&testM).Returning("*").Scan(t.Context())
	require.NoError(t, err)

	var quiz models.Quiz = models.Quiz{Path: fmt.Sprintf("/path/to/%v.md", rand.Text()), Checksum: []byte{}, Score: 1}
	err = db.NewInsert().Model(&quiz).Returning("*").Scan(t.Context())
	require.NoError(t, err)

	_, err = db.NewInsert().Model(&models.TestsQuizzies{TestUUID: testM.UUID, QuizPath: quiz.Path}).Exec(t.Context())
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

	testCases := []struct {
		desc      string
		testUUID  uuid.UUID
		userUUID  uuid.UUID
		groupUUID uuid.UUID
		quizPath  string
		expectErr bool
	}{
		{
			desc:      "Valid set",
			testUUID:  testM.UUID,
			quizPath:  quiz.Path,
			userUUID:  userUUID,
			groupUUID: groupUUID,
			expectErr: false,
		},
		{
			desc:      "Invalid test",
			testUUID:  uuid.Nil,
			quizPath:  quiz.Path,
			userUUID:  userUUID,
			groupUUID: groupUUID,
			expectErr: true,
		},
		{
			desc:      "Invalid quizPath",
			testUUID:  testM.UUID,
			quizPath:  "/unknown.md",
			userUUID:  userUUID,
			groupUUID: groupUUID,
			expectErr: true,
		},
		{
			desc:      "Invalid user",
			testUUID:  testM.UUID,
			quizPath:  quiz.Path,
			userUUID:  uuid.Nil,
			groupUUID: groupUUID,
			expectErr: true,
		},
		{
			desc:      "Invalid group",
			testUUID:  testM.UUID,
			quizPath:  quiz.Path,
			userUUID:  userUUID,
			groupUUID: uuid.Nil,
			expectErr: true,
		},
		{
			desc:      "Invalid group and user", // testing those pairs as they are only binded between eachother
			testUUID:  testM.UUID,
			quizPath:  quiz.Path,
			userUUID:  uuid.Nil,
			groupUUID: uuid.Nil,
			expectErr: true,
		},
		{
			desc:      "Invalid test and quiz",
			testUUID:  uuid.Nil,
			quizPath:  "/unknown.md",
			userUUID:  userUUID,
			groupUUID: groupUUID,
			expectErr: true,
		},
		{
			desc:      "Invalid all",
			testUUID:  uuid.Nil,
			quizPath:  "/unknown.md",
			userUUID:  uuid.Nil,
			groupUUID: uuid.Nil,
			expectErr: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Cleanup(func() {
				_, err := db.NewTruncateTable().Model((*models.Answer)(nil)).Exec(context.Background())
				require.NoError(t, err)
			})

			err := r.SetAnswer(t.Context(), tC.testUUID, tC.groupUUID, tC.userUUID, tC.quizPath, "", 1)
			if tC.expectErr {
				require.Error(t, err)
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

		err := r.SetAnswer(t.Context(), testM.UUID, groupUUID, userUUID, quiz.Path, "", 1)

		require.NoError(t, err)
		var insertion models.Answer = models.Answer{}
		err = db.NewSelect().Model(&insertion).
			Where("test_uuid = ?", testM.UUID).
			Where("group_uuid = ?", groupUUID).
			Where("user_uuid = ?", userUUID).
			Where("quiz_path = ?", quiz.Path).
			Scan(t.Context())
		require.NoError(t, err)

		err = r.SetAnswer(t.Context(), testM.UUID, groupUUID, userUUID, quiz.Path, "", 1)

		require.NoError(t, err)

		var updation models.Answer = models.Answer{}
		err = db.NewSelect().Model(&insertion).
			Where("test_uuid = ?", testM.UUID).
			Where("group_uuid = ?", groupUUID).
			Where("user_uuid = ?", userUUID).
			Where("quiz_path = ?", quiz.Path).
			Scan(t.Context())
		require.NoError(t, err)

		require.NotEqual(t, insertion.AnsweredAt, updation.AnsweredAt)
	})
}

func TestAnswer_Get(t *testing.T) {
	i := NewInjectorWithTestRepo(t)
	r := do.MustInvoke[answer.AnswerRepository](i)
	db := do.MustInvoke[*bun.DB](i)

	RegisterModels(db)

	testM := models.Test{Name: rand.Text()}
	err := db.NewInsert().Model(&testM).Returning("*").Scan(t.Context())
	require.NoError(t, err)

	var quiz models.Quiz = models.Quiz{Path: fmt.Sprintf("/path/to/%v.md", rand.Text()), Checksum: []byte{}, Score: 1}
	err = db.NewInsert().Model(&quiz).Returning("*").Scan(t.Context())
	require.NoError(t, err)

	_, err = db.NewInsert().Model(&models.TestsQuizzies{TestUUID: testM.UUID, QuizPath: quiz.Path}).Exec(t.Context())
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

	score := mrand.IntN(256)
	answerValue := rand.Text()
	_, err = db.NewInsert().Model(&models.Answer{
		GroupUUID:   groupUUID,
		TestUUID:    testM.UUID,
		UserUUID:    userUUID,
		QuizPath:    quiz.Path,
		AnswerValue: answerValue,
		Score:       score,
	}).Exec(t.Context())
	require.NoError(t, err)

	t.Run("Score suite", func(t *testing.T) {
		t.Run("Of known", func(t *testing.T) {
			scoreR, err := r.AnswerScore(t.Context(), userUUID, testM.UUID, groupUUID, quiz.Path)
			require.NoError(t, err)
			require.Equal(t, score, scoreR)
		})
		t.Run("Of unknown", func(t *testing.T) {
			scoreR, err := r.AnswerScore(t.Context(), uuid.Nil, uuid.Nil, uuid.Nil, "/unknown.md")
			require.Error(t, err)
			require.Equal(t, 0, scoreR)
		})
	})
	t.Run("Value suite", func(t *testing.T) {
		t.Run("Of known", func(t *testing.T) {
			answerR, err := r.Answer(t.Context(), userUUID, testM.UUID, groupUUID, quiz.Path)
			require.NoError(t, err)
			require.Equal(t, answerValue, answerR)
		})
		t.Run("Of unknown", func(t *testing.T) {
			answerR, err := r.Answer(t.Context(), uuid.Nil, uuid.Nil, uuid.Nil, "/unknown.md")
			require.Error(t, err)
			require.Equal(t, "", answerR)
		})
	})
}

func TestAnswer_Totalization(t *testing.T) {
	i := NewInjectorWithTestRepo(t)
	r := do.MustInvoke[answer.AnswerRepository](i)
	db := do.MustInvoke[*bun.DB](i)

	RegisterModels(db)

	testM := models.Test{Name: rand.Text()}
	err := db.NewInsert().Model(&testM).Returning("*").Scan(t.Context())
	require.NoError(t, err)

	var quiz models.Quiz = models.Quiz{Path: fmt.Sprintf("/path/to/%v.md", rand.Text()), Checksum: []byte{}, Score: 1}
	err = db.NewInsert().Model(&quiz).Returning("*").Scan(t.Context())
	require.NoError(t, err)

	_, err = db.NewInsert().Model(&models.TestsQuizzies{TestUUID: testM.UUID, QuizPath: quiz.Path}).Exec(t.Context())
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

	score := mrand.IntN(256)
	_, err = db.NewInsert().Model(&models.Answer{
		GroupUUID:   groupUUID,
		TestUUID:    testM.UUID,
		UserUUID:    userUUID,
		QuizPath:    quiz.Path,
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

		var scoreR int
		err = db.NewSelect().Model(&models.UserGroupsTests{
			UserUUID:  userUUID,
			GroupUUID: groupUUID,
			TestUUID:  testM.UUID,
		}).Column("score").Scan(t.Context(), &scoreR)

		require.NoError(t, err)
		require.Equal(t, score, scoreR)
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

		var scoreR int
		err = db.NewSelect().Model(&models.UserGroupsTests{
			UserUUID:  userUUID,
			GroupUUID: groupUUID,
			TestUUID:  testM.UUID,
		}).Column("score").Scan(t.Context(), &scoreR)

		require.NoError(t, err)
		require.Equal(t, score, scoreR)

		scoreN := mrand.IntN(256)
		var quiz models.Quiz = models.Quiz{Path: fmt.Sprintf("/path/to/%v.md", rand.Text()), Checksum: []byte{}, Score: 1}
		err = db.NewInsert().Model(&quiz).Returning("*").Scan(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.TestsQuizzies{TestUUID: testM.UUID, QuizPath: quiz.Path}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.Answer{
			GroupUUID:   groupUUID,
			TestUUID:    testM.UUID,
			UserUUID:    userUUID,
			QuizPath:    quiz.Path,
			AnswerValue: rand.Text(),
			Score:       scoreN,
		}).Exec(t.Context())
		require.NoError(t, err)

		err = r.Totalize(t.Context(), userUUID, testM.UUID, groupUUID)
		require.NoError(t, err)

		var scoreR2 int
		err = db.NewSelect().Model(&models.UserGroupsTests{
			UserUUID:  userUUID,
			GroupUUID: groupUUID,
			TestUUID:  testM.UUID,
		}).Column("score").Scan(t.Context(), &scoreR2)

		require.NoError(t, err)
		require.Equal(t, score+scoreN, scoreR2)
	})
}
