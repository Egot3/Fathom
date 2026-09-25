package handler_test

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	mrand "math/rand/v2"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/contracts"
	"github.com/egot3/fathom/internal/database/repositories"
	"github.com/egot3/fathom/internal/handler"
	"github.com/egot3/fathom/internal/models"
	"github.com/egot3/fathom/internal/quiz"
	testrunner "github.com/egot3/fathom/internal/testRunner"
	"github.com/egot3/fathom/internal/testutils"
	"github.com/egot3/fathom/server"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func randomTest(t *testing.T, db *bun.DB) uuid.UUID {
	t.Helper()

	var testUUID uuid.UUID
	err := db.NewInsert().Model(&models.Test{
		Name: rand.Text(),
	}).Returning("uuid").Scan(t.Context(), &testUUID)
	require.NoError(t, err)

	return testUUID
}

func randomQuiz(t *testing.T, db *bun.DB) (uuid.UUID, quiz.QuizAnswers) {
	t.Helper()

	fp, err := filepath.Abs("")
	require.NoError(t, err)

	answer := quiz.QuizAnswers{Input: new(quiz.AnswerInput{Input: rand.Text()})}
	b, err := json.Marshal(answer)
	require.NoError(t, err)

	var quizUUID uuid.UUID
	err = db.NewInsert().Model(&models.Quiz{
		Path:          filepath.VolumeName(fp) + string(filepath.Separator) + rand.Text() + ".md",
		Score:         1,
		Checksum:      [8]byte{},
		CorrectAnswer: string(b),
	}).Returning("uuid").Scan(t.Context(), &quizUUID)
	require.NoError(t, err)

	return quizUUID, answer
}

func predefinedQuiz(t *testing.T, db *bun.DB) (quizUUID uuid.UUID, path string) {
	t.Helper()

	f := testutils.TestQuiz(t)

	answer := quiz.QuizAnswers{Input: new(quiz.AnswerInput{Input: rand.Text()})}
	b, err := json.Marshal(answer)
	require.NoError(t, err)

	err = db.NewInsert().Model(&models.Quiz{
		Path:          f.Name(),
		Score:         1,
		Checksum:      [8]byte{},
		CorrectAnswer: string(b),
	}).Returning("uuid, path").Scan(t.Context(), &quizUUID, &path)
	require.NoError(t, err)

	return quizUUID, path
}

func randomGroup(t *testing.T, db *bun.DB) uuid.UUID {
	t.Helper()

	var groupUUID uuid.UUID
	err := db.NewInsert().Model(&models.Group{
		Name: rand.Text(),
	}).Returning("uuid").Scan(t.Context(), &groupUUID)
	require.NoError(t, err)

	return groupUUID
}

func randomUser(t *testing.T, db *bun.DB) uuid.UUID {
	t.Helper()

	var userUUID uuid.UUID
	err := db.NewInsert().Model(&models.User{
		Nickname:     rand.Text(),
		PasswordHash: []byte{},
		IsTeacher:    false,
	}).Returning("uuid").Scan(t.Context(), &userUUID)
	require.NoError(t, err)

	return userUUID
}

func TestTotalHandler_GetAnswer(t *testing.T) {

	t.Run("Valid", func(t *testing.T) {

		i := testutils.NewTestInjector(t,
			repositories.RepositoryPackage,
		)

		do.Provide(i, testrunner.NewManager)

		do.Provide(i, handler.NewTestService)
		router, err := server.ChiServer(i)
		require.NoError(t, err)

		db := do.MustInvoke[*bun.DB](i)

		userUUID := randomUser(t, db)
		groupUUID := randomGroup(t, db)

		testUUID := randomTest(t, db)
		quizUUID, correct := randomQuiz(t, db)

		_, err = db.NewInsert().Model(&models.TestsQuizzes{
			TestUUID: testUUID,
			QuizUUID: quizUUID,
			Position: 1,
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.GroupsUsers{
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
		}).Exec(t.Context())
		require.NoError(t, err)

		answer := quiz.QuizAnswers{Input: new(quiz.AnswerInput{Input: rand.Text()})}
		b, err := json.Marshal(answer)
		require.NoError(t, err)
		_, err = db.NewInsert().Model(&models.Answer{
			TestUUID:    testUUID,
			GroupUUID:   groupUUID,
			UserUUID:    userUUID,
			QuizUUID:    quizUUID,
			AnswerValue: string(b),
		}).Exec(t.Context())
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/total/%v/%v/%v/%v",
				groupUUID,
				userUUID,
				testUUID,
				quizUUID,
			),
			nil,
		)
		req.Header.Set("Content-Type", "application/json")
		testutils.AddTeacherCookie(t, req)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var bodyContract contracts.AnswerResponse
		err = json.NewDecoder(rec.Body).Decode(&bodyContract)
		require.NoError(t, err)

		require.Equal(t, answer, bodyContract.Answer.Chosen)
		require.Equal(t, correct, bodyContract.Answer.Correct)
	})

	t.Run("Not inferred", func(t *testing.T) {

		i := testutils.NewTestInjector(t,
			repositories.RepositoryPackage,
		)

		do.Provide(i, testrunner.NewManager)
		do.Provide(i, handler.NewTestService)

		db := do.MustInvoke[*bun.DB](i)

		userUUID := randomUser(t, db)
		groupUUID := randomGroup(t, db)

		testUUID := randomTest(t, db)
		quizUUID, _ := randomQuiz(t, db)

		_, err := db.NewInsert().Model(&models.TestsQuizzes{
			TestUUID: testUUID,
			QuizUUID: quizUUID,
			Position: 1,
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.GroupsUsers{
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
		}).Exec(t.Context())
		require.NoError(t, err)

		answer := rand.Text()
		_, err = db.NewInsert().Model(&models.Answer{
			TestUUID:    testUUID,
			GroupUUID:   groupUUID,
			UserUUID:    userUUID,
			QuizUUID:    quizUUID,
			AnswerValue: answer,
		}).Exec(t.Context())
		require.NoError(t, err)

		testCases := []struct {
			desc      string
			testUUID  uuid.UUID
			groupUUID uuid.UUID
			userUUID  uuid.UUID
			quizUUID  uuid.UUID
		}{
			{
				desc:      "No test",
				testUUID:  uuid.Max,
				groupUUID: groupUUID,
				userUUID:  userUUID,
				quizUUID:  quizUUID,
			},
			{
				desc:      "No quiz",
				testUUID:  testUUID,
				groupUUID: groupUUID,
				userUUID:  userUUID,
				quizUUID:  uuid.Max,
			},
			{
				desc:      "No quiz test pair",
				testUUID:  uuid.Max,
				groupUUID: groupUUID,
				userUUID:  userUUID,
				quizUUID:  uuid.Max,
			},
			{
				desc:      "No user",
				testUUID:  testUUID,
				groupUUID: groupUUID,
				userUUID:  uuid.Max,
				quizUUID:  quizUUID,
			},
			{
				desc:      "No group",
				testUUID:  testUUID,
				groupUUID: uuid.Max,
				userUUID:  userUUID,
				quizUUID:  quizUUID,
			},
			{
				desc:      "No user group pair",
				testUUID:  testUUID,
				groupUUID: uuid.Max,
				userUUID:  uuid.Max,
				quizUUID:  quizUUID,
			},
		}
		for _, tC := range testCases {
			t.Run(tC.desc, func(t *testing.T) {

				router, err := server.ChiServer(i)
				require.NoError(t, err)

				req := httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/total/%v/%v/%v/%v?page=0&size=5",
						tC.groupUUID,
						tC.userUUID,
						tC.testUUID,
						tC.quizUUID,
					),
					nil,
				)
				req.Header.Set("Content-Type", "application/json")
				testutils.AddTeacherCookie(t, req)
				rec := httptest.NewRecorder()

				router.ServeHTTP(rec, req)

				t.Log(rec.Body.String())

				var resp carefulness.JSONError
				err = json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)

				require.Equal(t, http.StatusNotFound, rec.Code, resp)
			})
		}
	})
}

func TestTotalHandler_GetGroupTotals(t *testing.T) {

	t.Run("Valid", func(t *testing.T) {

		i := testutils.NewTestInjector(t,
			repositories.RepositoryPackage,
		)

		do.Provide(i, testrunner.NewManager)

		do.Provide(i, handler.NewTestService)
		router, err := server.ChiServer(i)
		require.NoError(t, err)

		db := do.MustInvoke[*bun.DB](i)

		userUUID := randomUser(t, db)
		groupUUID := randomGroup(t, db)

		testUUID := randomTest(t, db)
		quizUUID, answer := randomQuiz(t, db)

		_, err = db.NewInsert().Model(&models.TestsQuizzes{
			TestUUID: testUUID,
			QuizUUID: quizUUID,
			Position: 1,
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.GroupsUsers{
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
		}).Exec(t.Context())
		require.NoError(t, err)

		b, err := json.Marshal(answer)
		require.NoError(t, err)
		_, err = db.NewInsert().Model(&models.Answer{
			TestUUID:    testUUID,
			GroupUUID:   groupUUID,
			UserUUID:    userUUID,
			QuizUUID:    quizUUID,
			AnswerValue: string(b),
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.UserGroupsTests{
			TestUUID:  testUUID,
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
			Score:     1,
		}).Exec(t.Context())
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/total/%v/all/%v?page=0&size=5",
				groupUUID,
				testUUID,
			),
			nil,
		)
		req.Header.Set("Content-Type", "application/json")
		testutils.AddTeacherCookie(t, req)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code, "there was an early body consumtion")

		var bodyContract contracts.TotalsResponse
		err = json.NewDecoder(rec.Body).Decode(&bodyContract)
		require.NoError(t, err)

		require.Len(t, bodyContract.Totals, 1)

		require.EqualValues(t, 1, bodyContract.Totals[0].Score)
	})

	t.Run("Not inferred", func(t *testing.T) {

		i := testutils.NewTestInjector(t,
			repositories.RepositoryPackage,
		)

		do.Provide(i, testrunner.NewManager)
		do.Provide(i, handler.NewTestService)

		db := do.MustInvoke[*bun.DB](i)

		userUUID := randomUser(t, db)
		groupUUID := randomGroup(t, db)

		testUUID := randomTest(t, db)
		quizUUID, _ := randomQuiz(t, db)

		_, err := db.NewInsert().Model(&models.TestsQuizzes{
			TestUUID: testUUID,
			QuizUUID: quizUUID,
			Position: 1,
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.GroupsUsers{
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
		}).Exec(t.Context())
		require.NoError(t, err)

		answer := rand.Text()
		_, err = db.NewInsert().Model(&models.Answer{
			TestUUID:    testUUID,
			GroupUUID:   groupUUID,
			UserUUID:    userUUID,
			QuizUUID:    quizUUID,
			AnswerValue: answer,
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.UserGroupsTests{
			TestUUID:  testUUID,
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
			Score:     1,
		}).Exec(t.Context())
		require.NoError(t, err)

		testCases := []struct {
			desc      string
			testUUID  uuid.UUID
			groupUUID uuid.UUID
		}{
			{
				desc:      "No test",
				testUUID:  uuid.Max,
				groupUUID: groupUUID,
			},
			{
				desc:      "No group",
				testUUID:  testUUID,
				groupUUID: uuid.Max,
			},
			{
				desc:      "Nothing.",
				testUUID:  uuid.Max,
				groupUUID: uuid.Max,
			},
		}
		for _, tC := range testCases {
			t.Run(tC.desc, func(t *testing.T) {

				router, err := server.ChiServer(i)
				require.NoError(t, err)

				req := httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/total/%v/all/%v?page=0&size=5",
						tC.groupUUID,
						tC.testUUID,
					),
					nil,
				)
				req.Header.Set("Content-Type", "application/json")
				testutils.AddTeacherCookie(t, req)
				rec := httptest.NewRecorder()

				router.ServeHTTP(rec, req)

				t.Log(rec.Body.String())

				var resp carefulness.JSONError
				err = json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)

				require.Equal(t, http.StatusNotFound, rec.Code, resp)
			})
		}
	})
}

func TestTotalHandler_GetTestTotals(t *testing.T) {

	t.Run("Valid", func(t *testing.T) {

		i := testutils.NewTestInjector(t,
			repositories.RepositoryPackage,
		)

		do.Provide(i, testrunner.NewManager)

		do.Provide(i, handler.NewTestService)
		router, err := server.ChiServer(i)
		require.NoError(t, err)

		db := do.MustInvoke[*bun.DB](i)

		userUUID := randomUser(t, db)
		groupUUID := randomGroup(t, db)

		testUUID := randomTest(t, db)
		quizUUID, answer := randomQuiz(t, db)

		_, err = db.NewInsert().Model(&models.TestsQuizzes{
			TestUUID: testUUID,
			QuizUUID: quizUUID,
			Position: 1,
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.GroupsUsers{
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
		}).Exec(t.Context())
		require.NoError(t, err)

		b, err := json.Marshal(answer)
		require.NoError(t, err)
		_, err = db.NewInsert().Model(&models.Answer{
			TestUUID:    testUUID,
			GroupUUID:   groupUUID,
			UserUUID:    userUUID,
			QuizUUID:    quizUUID,
			AnswerValue: string(b),
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.UserGroupsTests{
			TestUUID:  testUUID,
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
			Score:     1,
		}).Exec(t.Context())
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/total/all/all/%v?page=0&size=5",
				testUUID,
			),
			nil,
		)
		req.Header.Set("Content-Type", "application/json")
		testutils.AddTeacherCookie(t, req)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code, "there was an early body consumtion")

		var bodyContract contracts.TotalsResponse
		err = json.NewDecoder(rec.Body).Decode(&bodyContract)
		require.NoError(t, err)

		require.Len(t, bodyContract.Totals, 1)

		require.EqualValues(t, 1, bodyContract.Totals[0].Score)
	})

	t.Run("Not inferred", func(t *testing.T) {

		i := testutils.NewTestInjector(t,
			repositories.RepositoryPackage,
		)

		do.Provide(i, testrunner.NewManager)
		do.Provide(i, handler.NewTestService)

		db := do.MustInvoke[*bun.DB](i)

		userUUID := randomUser(t, db)
		groupUUID := randomGroup(t, db)

		testUUID := randomTest(t, db)
		quizUUID, _ := randomQuiz(t, db)

		_, err := db.NewInsert().Model(&models.TestsQuizzes{
			TestUUID: testUUID,
			QuizUUID: quizUUID,
			Position: 1,
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.GroupsUsers{
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
		}).Exec(t.Context())
		require.NoError(t, err)

		answer := rand.Text()
		_, err = db.NewInsert().Model(&models.Answer{
			TestUUID:    testUUID,
			GroupUUID:   groupUUID,
			UserUUID:    userUUID,
			QuizUUID:    quizUUID,
			AnswerValue: answer,
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.UserGroupsTests{
			TestUUID:  testUUID,
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
			Score:     1,
		}).Exec(t.Context())
		require.NoError(t, err)

		testCases := []struct {
			desc     string
			testUUID uuid.UUID
		}{
			{
				desc:     "No test",
				testUUID: uuid.Max,
			},
		}
		for _, tC := range testCases {
			t.Run(tC.desc, func(t *testing.T) {

				router, err := server.ChiServer(i)
				require.NoError(t, err)

				req := httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/total/all/all/%v?page=0&size=5",
						tC.testUUID,
					),
					nil,
				)
				req.Header.Set("Content-Type", "application/json")
				testutils.AddTeacherCookie(t, req)
				rec := httptest.NewRecorder()

				router.ServeHTTP(rec, req)

				t.Log(rec.Body.String())

				var resp carefulness.JSONError
				err = json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)

				require.Equal(t, http.StatusNotFound, rec.Code, resp)
			})
		}
	})
}

func TestTotalHandler_GetUserTotal(t *testing.T) {

	t.Run("Valid", func(t *testing.T) {

		i := testutils.NewTestInjector(t,
			repositories.RepositoryPackage,
		)

		do.Provide(i, testrunner.NewManager)

		do.Provide(i, handler.NewTestService)
		router, err := server.ChiServer(i)
		require.NoError(t, err)

		db := do.MustInvoke[*bun.DB](i)

		userUUID := randomUser(t, db)
		groupUUID := randomGroup(t, db)

		testUUID := randomTest(t, db)
		quizUUID, answer := randomQuiz(t, db)

		_, err = db.NewInsert().Model(&models.TestsQuizzes{
			TestUUID: testUUID,
			QuizUUID: quizUUID,
			Position: 1,
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.GroupsUsers{
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
		}).Exec(t.Context())
		require.NoError(t, err)

		b, err := json.Marshal(answer)
		require.NoError(t, err)
		_, err = db.NewInsert().Model(&models.Answer{
			TestUUID:    testUUID,
			GroupUUID:   groupUUID,
			UserUUID:    userUUID,
			QuizUUID:    quizUUID,
			AnswerValue: string(b),
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.UserGroupsTests{
			TestUUID:  testUUID,
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
			Score:     1,
		}).Exec(t.Context())
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/total/%v/%v/%v?page=0&size=1",
				groupUUID,
				userUUID,
				testUUID,
			),
			nil,
		)
		req.Header.Set("Content-Type", "application/json")
		testutils.AddTeacherCookie(t, req)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		body := rec.Body.Bytes()

		require.Equal(t, http.StatusOK, rec.Code, string(body))

		var bodyContract contracts.TotalsResponse
		err = json.Unmarshal(body, &bodyContract)
		require.NoError(t, err, string(body))

		require.Len(t, bodyContract.Totals, 1)
		require.EqualValues(t, 1, bodyContract.Totals[0].Score)
	})

	t.Run("Not inferred", func(t *testing.T) {

		i := testutils.NewTestInjector(t,
			repositories.RepositoryPackage,
		)

		do.Provide(i, testrunner.NewManager)
		do.Provide(i, handler.NewTestService)

		db := do.MustInvoke[*bun.DB](i)

		userUUID := randomUser(t, db)
		groupUUID := randomGroup(t, db)

		testUUID := randomTest(t, db)
		quizUUID, _ := randomQuiz(t, db)

		_, err := db.NewInsert().Model(&models.TestsQuizzes{
			TestUUID: testUUID,
			QuizUUID: quizUUID,
			Position: 1,
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.GroupsUsers{
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
		}).Exec(t.Context())
		require.NoError(t, err)

		answer := rand.Text()
		_, err = db.NewInsert().Model(&models.Answer{
			TestUUID:    testUUID,
			GroupUUID:   groupUUID,
			UserUUID:    userUUID,
			QuizUUID:    quizUUID,
			AnswerValue: answer,
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.UserGroupsTests{
			TestUUID:  testUUID,
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
			Score:     1,
		}).Exec(t.Context())
		require.NoError(t, err)

		testCases := []struct {
			desc      string
			testUUID  uuid.UUID
			groupUUID uuid.UUID
			userUUID  uuid.UUID
		}{
			{
				desc:      "No test",
				testUUID:  uuid.Max,
				groupUUID: groupUUID,
				userUUID:  userUUID,
			},
			{
				desc:      "No group",
				testUUID:  testUUID,
				groupUUID: uuid.Max,
				userUUID:  userUUID,
			},
			{
				desc:      "No user",
				testUUID:  testUUID,
				groupUUID: groupUUID,
				userUUID:  uuid.Max,
			},
			{
				desc:      "No group and user pair",
				testUUID:  testUUID,
				groupUUID: uuid.Max,
				userUUID:  uuid.Max,
			},
			{
				desc:      "Nothing.",
				testUUID:  uuid.Max,
				groupUUID: uuid.Max,
				userUUID:  uuid.Max,
			},
		}
		for _, tC := range testCases {
			t.Run(tC.desc, func(t *testing.T) {

				router, err := server.ChiServer(i)
				require.NoError(t, err)

				req := httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/total/%v/%v/%v?page=0&size=1",
						tC.groupUUID,
						tC.userUUID,
						tC.testUUID,
					),
					nil,
				)
				req.Header.Set("Content-Type", "application/json")
				testutils.AddTeacherCookie(t, req)
				rec := httptest.NewRecorder()

				router.ServeHTTP(rec, req)

				t.Log(rec.Body.String())

				var resp carefulness.JSONError
				err = json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)

				require.Equal(t, http.StatusNotFound, rec.Code, resp)
			})
		}
	})
}

func TestTotalHandler_GetUserTotals(t *testing.T) {

	t.Run("Valid", func(t *testing.T) {

		i := testutils.NewTestInjector(t,
			repositories.RepositoryPackage,
		)

		do.Provide(i, testrunner.NewManager)

		do.Provide(i, handler.NewTestService)
		router, err := server.ChiServer(i)
		require.NoError(t, err)

		db := do.MustInvoke[*bun.DB](i)

		userUUID := randomUser(t, db)
		groupUUID := randomGroup(t, db)

		testUUID := randomTest(t, db)
		quizUUID, answer := randomQuiz(t, db)

		_, err = db.NewInsert().Model(&models.TestsQuizzes{
			TestUUID: testUUID,
			QuizUUID: quizUUID,
			Position: 1,
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.GroupsUsers{
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
		}).Exec(t.Context())
		require.NoError(t, err)

		b, err := json.Marshal(answer)
		require.NoError(t, err)
		_, err = db.NewInsert().Model(&models.Answer{
			TestUUID:    testUUID,
			GroupUUID:   groupUUID,
			UserUUID:    userUUID,
			QuizUUID:    quizUUID,
			AnswerValue: string(b),
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.UserGroupsTests{
			TestUUID:  testUUID,
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
			Score:     1,
		}).Exec(t.Context())
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/total/all/%v/all?page=0&size=1",
				userUUID,
			),
			nil,
		)
		req.Header.Set("Content-Type", "application/json")
		testutils.AddTeacherCookie(t, req)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		body := rec.Body.Bytes()

		require.Equal(t, http.StatusOK, rec.Code, string(body))

		var bodyContract contracts.TotalsResponse
		err = json.Unmarshal(body, &bodyContract)
		require.NoError(t, err)

		require.Len(t, bodyContract.Totals, 1)

		require.EqualValues(t, 1, bodyContract.Totals[0].Score)
	})

	t.Run("Not inferred", func(t *testing.T) {

		i := testutils.NewTestInjector(t,
			repositories.RepositoryPackage,
		)

		do.Provide(i, testrunner.NewManager)
		do.Provide(i, handler.NewTestService)

		db := do.MustInvoke[*bun.DB](i)

		userUUID := randomUser(t, db)
		groupUUID := randomGroup(t, db)

		testUUID := randomTest(t, db)
		quizUUID, _ := randomQuiz(t, db)

		_, err := db.NewInsert().Model(&models.TestsQuizzes{
			TestUUID: testUUID,
			QuizUUID: quizUUID,
			Position: 1,
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.GroupsUsers{
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
		}).Exec(t.Context())
		require.NoError(t, err)

		answer := rand.Text()
		_, err = db.NewInsert().Model(&models.Answer{
			TestUUID:    testUUID,
			GroupUUID:   groupUUID,
			UserUUID:    userUUID,
			QuizUUID:    quizUUID,
			AnswerValue: answer,
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.UserGroupsTests{
			TestUUID:  testUUID,
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
			Score:     1,
		}).Exec(t.Context())
		require.NoError(t, err)

		t.Run("No user", func(t *testing.T) {

			router, err := server.ChiServer(i)
			require.NoError(t, err)

			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/api/v1/total/all/%v/all?page=0&size=1",
					uuid.Max,
				),
				nil,
			)
			req.Header.Set("Content-Type", "application/json")
			testutils.AddTeacherCookie(t, req)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			var resp carefulness.JSONError
			err = json.NewDecoder(rec.Body).Decode(&resp)
			require.NoError(t, err)

			require.Equal(t, http.StatusNotFound, rec.Code, resp)
		})

	})
}

func TestTotalHandler_PostAnswer(t *testing.T) {

	t.Run("Valid", func(t *testing.T) {
		i := testutils.NewTestInjector(t, repositories.RepositoryPackage)

		do.Provide(i, testrunner.NewManager)
		do.Provide(i, handler.NewTestService)

		m := do.MustInvoke[*testrunner.Manager](i)

		db := do.MustInvoke[*bun.DB](i)

		userUUID := randomUser(t, db)
		groupUUID := randomGroup(t, db)

		testUUID := randomTest(t, db)
		quizUUID, path := predefinedQuiz(t, db)

		_, err := db.NewInsert().Model(&models.TestsQuizzes{
			TestUUID: testUUID,
			QuizUUID: quizUUID,
			Position: 1,
		}).Exec(t.Context())
		require.NoError(t, err)

		tr, err := m.Start(t.Context(), 30*time.Second, []string{path}, uuid.UUIDs{quizUUID}, uuid.UUIDs{groupUUID}, testUUID)
		require.NoError(t, err)

		all := m.GetAll()
		require.Len(t, all, 1)
		key := all[0]
		trc, ok := m.Get(key)
		require.True(t, ok)
		require.Equal(t, trc, tr)

		_, err = db.NewInsert().Model(&models.GroupsUsers{
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
		}).Exec(t.Context())
		require.NoError(t, err)

		router, err := server.ChiServer(i)
		require.NoError(t, err)

		reqJSON, err := json.Marshal(contracts.PostAnswerRequest{
			Value: quiz.QuizAnswers{
				Input: &quiz.AnswerInput{
					Input: "1",
				},
			},
		})
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodPost,
			fmt.Sprintf("/api/v1/total/%v/%v/running/%d/%v",
				groupUUID, userUUID,
				key, quizUUID,
			),
			bytes.NewReader(reqJSON),
		)
		req.Header.Set("Content-Type", "application/json")
		testutils.AddTeacherCookie(t, req)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNoContent, rec.Code, rec.Body.String())

		var answer models.Answer
		err = db.NewSelect().Model(&answer).
			Where("user_uuid = ?", userUUID).
			Where("group_uuid = ?", groupUUID).
			Where("test_uuid = ?", testUUID).
			Where("quiz_uuid = ?", quizUUID).Scan(t.Context())
		require.Equal(t, "{\"input\":{\"input\":\"1\"}}", answer.AnswerValue) // да, хранить объекты в json-е иногда нормально
	})

	t.Run("Not running", func(t *testing.T) {
		i := testutils.NewTestInjector(t, repositories.RepositoryPackage)

		do.Provide(i, testrunner.NewManager)
		do.Provide(i, handler.NewTestService)

		db := do.MustInvoke[*bun.DB](i)

		userUUID := randomUser(t, db)
		groupUUID := randomGroup(t, db)

		testUUID := randomTest(t, db)
		quizUUID, _ := predefinedQuiz(t, db)

		_, err := db.NewInsert().Model(&models.TestsQuizzes{
			TestUUID: testUUID,
			QuizUUID: quizUUID,
			Position: 1,
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.GroupsUsers{
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
		}).Exec(t.Context())
		require.NoError(t, err)

		router, err := server.ChiServer(i)
		require.NoError(t, err)

		reqJSON, err := json.Marshal(contracts.PostAnswerRequest{
			Value: quiz.QuizAnswers{
				Input: &quiz.AnswerInput{
					Input: "1",
				},
			},
		})
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodPost,
			fmt.Sprintf("/api/v1/total/%v/%v/running/%d/%v",
				groupUUID, userUUID,
				mrand.Uint64(), quizUUID,
			),
			bytes.NewReader(reqJSON),
		)
		req.Header.Set("Content-Type", "application/json")
		testutils.AddTeacherCookie(t, req)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusLocked, rec.Code, rec.Body.String())
	})

	t.Run("Not found", func(t *testing.T) {
		i := testutils.NewTestInjector(t, repositories.RepositoryPackage)

		do.Provide(i, testrunner.NewManager)
		do.Provide(i, handler.NewTestService)

		m := do.MustInvoke[*testrunner.Manager](i)

		db := do.MustInvoke[*bun.DB](i)

		userUUID := randomUser(t, db)
		groupUUID := randomGroup(t, db)

		testUUID := randomTest(t, db)
		quizUUID, path := predefinedQuiz(t, db)

		_, err := db.NewInsert().Model(&models.TestsQuizzes{
			TestUUID: testUUID,
			QuizUUID: quizUUID,
			Position: 1,
		}).Exec(t.Context())
		require.NoError(t, err)

		tr, err := m.Start(t.Context(), 30*time.Second, []string{path}, uuid.UUIDs{quizUUID}, uuid.UUIDs{groupUUID}, testUUID)
		require.NoError(t, err)

		all := m.GetAll()
		require.Len(t, all, 1)
		key := all[0]
		trc, ok := m.Get(key)
		require.True(t, ok)
		require.Equal(t, trc, tr)

		_, err = db.NewInsert().Model(&models.GroupsUsers{
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
		}).Exec(t.Context())
		require.NoError(t, err)

		router, err := server.ChiServer(i)
		require.NoError(t, err)

		testCases := []struct {
			desc      string
			groupUUID uuid.UUID
			userUUID  uuid.UUID
			quizUUID  uuid.UUID
		}{
			{
				desc:      "No group",
				groupUUID: uuid.Max,
				userUUID:  userUUID,
				quizUUID:  quizUUID,
			},
			{
				desc:      "No user",
				groupUUID: groupUUID,
				userUUID:  uuid.Max,
				quizUUID:  quizUUID,
			},
			{
				desc:      "No group and user",
				groupUUID: uuid.Max,
				userUUID:  uuid.Max,
				quizUUID:  quizUUID,
			},
			{
				desc:      "No quiz",
				groupUUID: groupUUID,
				userUUID:  userUUID,
				quizUUID:  uuid.Max,
			},
		}
		for _, tC := range testCases {
			t.Run(tC.desc, func(t *testing.T) {
				reqJSON, err := json.Marshal(contracts.PostAnswerRequest{
					Value: quiz.QuizAnswers{
						Input: &quiz.AnswerInput{
							Input: "1",
						},
					},
				})
				require.NoError(t, err)

				req := httptest.NewRequest(
					http.MethodPost,
					fmt.Sprintf("/api/v1/total/%v/%v/running/%d/%v",
						tC.groupUUID, tC.userUUID,
						key, tC.quizUUID,
					),
					bytes.NewReader(reqJSON),
				)
				req.Header.Set("Content-Type", "application/json")
				testutils.AddTeacherCookie(t, req)
				rec := httptest.NewRecorder()

				router.ServeHTTP(rec, req)

				require.Equal(t, http.StatusNotFound, rec.Code, rec.Body.String())
			})
		}
	})
}

func TestTotalHandler_Totalize(t *testing.T) {

	t.Run("Valid", func(t *testing.T) {
		i := testutils.NewTestInjector(t, repositories.RepositoryPackage)

		do.Provide(i, testrunner.NewManager)
		do.Provide(i, handler.NewTestService)

		db := do.MustInvoke[*bun.DB](i)

		userUUID := randomUser(t, db)
		groupUUID := randomGroup(t, db)

		testUUID := randomTest(t, db)
		quizUUID, path := predefinedQuiz(t, db)

		t.Logf("userUUID: %v, groupUUID: %v, quizUUID: %v, testUUID: %v", userUUID, groupUUID, quizUUID, testUUID)

		_, err := db.NewInsert().Model(&models.TestsQuizzes{
			TestUUID: testUUID,
			QuizUUID: quizUUID,
			Position: 1,
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.GroupsUsers{
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.Answer{
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
			TestUUID:  testUUID,
			QuizUUID:  quizUUID,
			Score:     1,
		}).Exec(t.Context())
		require.NoError(t, err)

		router, err := server.ChiServer(i)
		require.NoError(t, err)

		m := do.MustInvoke[*testrunner.Manager](i)

		tr, err := m.Start(t.Context(), 30*time.Second, []string{path}, uuid.UUIDs{quizUUID}, uuid.UUIDs{groupUUID}, testUUID)
		require.NoError(t, err)

		all := m.GetAll()
		require.Len(t, all, 1)
		key := all[0]
		trc, ok := m.Get(key)
		require.True(t, ok)
		require.Equal(t, trc, tr)

		req := httptest.NewRequest(
			http.MethodPost,
			fmt.Sprintf("/api/v1/total/%v/%v/running/%d",
				groupUUID, userUUID, key,
			),
			nil,
		)
		req.Header.Set("Content-Type", "application/json")
		testutils.AddTeacherCookie(t, req)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNoContent, rec.Code, rec.Body.String())

		var total models.UserGroupsTests
		err = db.NewSelect().Model(&total).
			Where("user_uuid = ?", userUUID).
			Where("group_uuid = ?", groupUUID).
			Where("test_uuid = ?", testUUID).Scan(t.Context())
		require.NoError(t, err)
		require.EqualValues(t, 1, total.Score, fmt.Sprintf("%+v", total))
	})

	t.Run("Not running", func(t *testing.T) {
		i := testutils.NewTestInjector(t, repositories.RepositoryPackage)

		do.Provide(i, testrunner.NewManager)
		do.Provide(i, handler.NewTestService)

		db := do.MustInvoke[*bun.DB](i)

		userUUID := randomUser(t, db)
		groupUUID := randomGroup(t, db)

		testUUID := randomTest(t, db)
		quizUUID, _ := predefinedQuiz(t, db)

		_, err := db.NewInsert().Model(&models.TestsQuizzes{
			TestUUID: testUUID,
			QuizUUID: quizUUID,
			Position: 1,
		}).Exec(t.Context())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.GroupsUsers{
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
		}).Exec(t.Context())
		require.NoError(t, err)

		router, err := server.ChiServer(i)
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodPost,
			fmt.Sprintf("/api/v1/total/%v/%v/running/%d",
				groupUUID, userUUID, mrand.Uint64(),
			),
			nil,
		)
		req.Header.Set("Content-Type", "application/json")
		testutils.AddTeacherCookie(t, req)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusLocked, rec.Code, rec.Body.String())
	})

	t.Run("Not found", func(t *testing.T) {
		i := testutils.NewTestInjector(t, repositories.RepositoryPackage)

		do.Provide(i, testrunner.NewManager)
		do.Provide(i, handler.NewTestService)

		db := do.MustInvoke[*bun.DB](i)

		userUUID := randomUser(t, db)
		groupUUID := randomGroup(t, db)

		testUUID := randomTest(t, db)
		quizUUID, path := predefinedQuiz(t, db)

		_, err := db.NewInsert().Model(&models.TestsQuizzes{
			TestUUID: testUUID,
			QuizUUID: quizUUID,
			Position: 1,
		}).Exec(t.Context())
		require.NoError(t, err)

		require.NoError(t, err)

		_, err = db.NewInsert().Model(&models.GroupsUsers{
			GroupUUID: groupUUID,
			UserUUID:  userUUID,
		}).Exec(t.Context())
		require.NoError(t, err)

		router, err := server.ChiServer(i)
		require.NoError(t, err)

		testCases := []struct {
			desc      string
			groupUUID uuid.UUID
			userUUID  uuid.UUID
		}{
			{
				desc:      "No group",
				groupUUID: uuid.Max,
				userUUID:  userUUID,
			},
			{
				desc:      "No user",
				groupUUID: groupUUID,
				userUUID:  uuid.Max,
			},
			{
				desc:      "No group and user",
				groupUUID: uuid.Max,
				userUUID:  uuid.Max,
			},
		}
		for _, tC := range testCases {
			t.Run(tC.desc, func(t *testing.T) {
				reqJSON, err := json.Marshal(contracts.PostAnswerRequest{
					Value: quiz.QuizAnswers{
						Input: &quiz.AnswerInput{
							Input: "1",
						},
					},
				})
				require.NoError(t, err)

				i := i.Scope(tC.desc)
				do.Provide(i, testrunner.NewManager)

				m := do.MustInvoke[*testrunner.Manager](i)

				tr, err := m.Start(t.Context(), 30*time.Second, []string{path}, uuid.UUIDs{quizUUID}, uuid.UUIDs{groupUUID}, testUUID)
				require.NoError(t, err)

				all := m.GetAll()
				require.Len(t, all, 1)
				key := all[0]
				trc, ok := m.Get(key)
				require.True(t, ok)
				require.Equal(t, trc, tr)

				req := httptest.NewRequest(
					http.MethodPost,
					fmt.Sprintf("/api/v1/total/%v/%v/running/%d",
						tC.groupUUID, tC.userUUID, key,
					),
					bytes.NewReader(reqJSON),
				)
				req.Header.Set("Content-Type", "application/json")
				testutils.AddTeacherCookie(t, req)
				rec := httptest.NewRecorder()

				router.ServeHTTP(rec, req)

				require.Equal(t, http.StatusLocked, rec.Code, rec.Body.String())
			})
		}
	})
}
