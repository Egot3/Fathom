package handler_test

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	mrand "math/rand/v2"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	jwtutils "github.com/egot3/fathom/internal/JWTutils"
	"github.com/egot3/fathom/internal/config"
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

func TestQuizHandler_Post(t *testing.T) {
	t.Parallel()

	token, err := jwtutils.GenerateToken(uuid.Nil, true)
	require.NoError(t, err)

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()

		i := testutils.NewTestInjector(t,
			repositories.RepositoryPackage,
		)
		do.ProvideValue(i, slog.Default())
		do.Provide(i, testrunner.NewTestRunner)

		do.Provide(i, handler.NewTestService)
		router, err := server.ChiServer(i)
		require.NoError(t, err)

		name := rand.Text()
		body := `# quiz!
			 there is a body!
			 [yeah!]`
		reqJSON, _ := json.Marshal(contracts.PostQuizRequest{
			Name: name,
			Body: body,
			Meta: quiz.Frontmatter{
				Kind:  quiz.Input,
				Score: 1,
			},
		})
		req := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/quiz/",
			bytes.NewReader(reqJSON),
		)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		bodyString := rec.Body.String()
		require.Equal(t, http.StatusNoContent, rec.Code, bodyString)

		path := config.TurnToAbs(name)
		t.Logf("Filepath: %v", path)
		require.FileExists(t, path)
		b, err := os.ReadFile(path)
		require.NoError(t, err)

		require.Contains(t, string(b), body, string(b))

		err = os.Remove(path)
		require.NoError(t, err)
	})

	t.Run("Orphans", func(t *testing.T) {
		meta := quiz.Frontmatter{
			Kind:  quiz.Input,
			Score: 1,
		}
		testCases := []struct {
			desc string
			name string
			body string
			meta quiz.Frontmatter
		}{
			{
				desc: "No name",
				name: "",
				body: rand.Text(),
				meta: meta,
			},
			{
				desc: "No body",
				name: rand.Text(),
				body: "",
				meta: meta,
			},
			{
				desc: "No meta",
				name: rand.Text(),
				body: rand.Text(),
				meta: quiz.Frontmatter{},
			},
			{
				desc: "Nothing",
				name: "",
				body: "",
				meta: quiz.Frontmatter{},
			},
		}
		for _, tC := range testCases {
			t.Run(tC.desc, func(t *testing.T) {
				t.Parallel()

				i := testutils.NewTestInjector(t,
					repositories.RepositoryPackage,
				)
				do.ProvideValue(i, slog.Default())
				do.Provide(i, testrunner.NewTestRunner)

				do.Provide(i, handler.NewTestService)
				router, err := server.ChiServer(i)
				require.NoError(t, err)

				reqJSON, _ := json.Marshal(contracts.PostQuizRequest{
					Name: tC.name,
					Body: tC.body,
					Meta: tC.meta,
				})
				req := httptest.NewRequest(
					http.MethodPost,
					"/api/v1/quiz/",
					bytes.NewReader(reqJSON),
				)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

				rec := httptest.NewRecorder()

				router.ServeHTTP(rec, req)

				bodyString := rec.Body.String()
				require.Equal(t, http.StatusBadRequest, rec.Code, bodyString)
			})
		}
	})

	t.Run("Conflict!", func(t *testing.T) {
		t.Parallel()

		i := testutils.NewTestInjector(t,
			repositories.RepositoryPackage,
		)
		do.ProvideValue(i, slog.Default())
		do.Provide(i, testrunner.NewTestRunner)

		name := rand.Text()
		db := do.MustInvoke[*bun.DB](i)
		_, err := db.NewInsert().Model(&models.Quiz{
			Path:          config.TurnToAbs(name),
			Checksum:      []byte{},
			Score:         1,
			CorrectAnswer: "omega",
		}).Exec(t.Context())

		do.Provide(i, handler.NewTestService)
		router, err := server.ChiServer(i)
		require.NoError(t, err)

		body := `# quiz!
			 there is a body!
			 [yeah!]`
		reqJSON, _ := json.Marshal(contracts.PostQuizRequest{
			Name: name,
			Body: body,
			Meta: quiz.Frontmatter{
				Kind:  quiz.Input,
				Score: 1,
			},
		})
		req := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/quiz/",
			bytes.NewReader(reqJSON),
		)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		bodyString := rec.Body.String()
		require.Equal(t, http.StatusConflict, rec.Code, bodyString)

		path := config.TurnToAbs(name)
		require.NoFileExists(t, path)
	})
}

func TestQuizHandler_Get(t *testing.T) {
	t.Parallel()

	token, err := jwtutils.GenerateToken(uuid.Nil, true)
	require.NoError(t, err)

	i := testutils.NewTestInjector(t,
		repositories.RepositoryPackage,
	)
	do.ProvideValue(i, slog.Default())
	do.Provide(i, testrunner.NewTestRunner)

	var quizUUID uuid.UUID
	score := 1
	db := do.MustInvoke[*bun.DB](i)
	err = db.NewInsert().Model(&models.Quiz{
		Path:          "/home/ETS/programming/Fathom/application/internal/testutils/placebo.md",
		Checksum:      []byte{},
		Score:         score,
		CorrectAnswer: "x",
	}).Returning("uuid").
		Scan(t.Context(), &quizUUID)
	require.NoError(t, err)

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()

		i = i.Scope("valid")

		do.Provide(i, handler.NewTestService)
		router, err := server.ChiServer(i)
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/quiz/%v", quizUUID),
			nil,
		)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		bodyString := rec.Body.String()
		require.Equal(t, http.StatusOK, rec.Code, bodyString)

		var quizS contracts.GetQuizResponse
		err = json.Unmarshal(rec.Body.Bytes(), &quizS)
		require.NoError(t, err)

		require.Equal(t, score, quizS.Quiz.Meta.Score)
		require.Equal(t, "x", quizS.Quiz.Answer.Input.Input)
	})

	t.Run("Invalid", func(t *testing.T) {
		t.Parallel()

		i = i.Scope("invalid")

		do.Provide(i, handler.NewTestService)
		router, err := server.ChiServer(i)
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/quiz/%v", uuid.Max),
			nil,
		)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		bodyString := rec.Body.String()
		require.Equal(t, http.StatusNotFound, rec.Code, bodyString)
	})
}

func TestQuizHandler_List(t *testing.T) {
	t.Parallel()

	token, err := jwtutils.GenerateToken(uuid.Nil, true)
	require.NoError(t, err)

	i := testutils.NewTestInjector(t,
		repositories.RepositoryPackage,
	)
	do.ProvideValue(i, slog.Default())
	do.Provide(i, testrunner.NewTestRunner)

	var quizUUID uuid.UUID
	score := 1
	db := do.MustInvoke[*bun.DB](i)
	err = db.NewInsert().Model(&models.Quiz{
		Path:          "/home/ETS/programming/Fathom/application/internal/testutils/placebo.md",
		Checksum:      []byte{},
		Score:         score,
		CorrectAnswer: "x",
	}).Returning("uuid").
		Scan(t.Context(), &quizUUID)
	require.NoError(t, err)

	do.Provide(i, handler.NewTestService)
	router, err := server.ChiServer(i)
	require.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/quiz?page=0&size=1",
		nil,
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	bodyString := rec.Body.String()
	require.Equal(t, http.StatusOK, rec.Code, bodyString)

	var quizL contracts.ListQuizResponse
	err = json.Unmarshal(rec.Body.Bytes(), &quizL)
	require.NoError(t, err)

	require.Len(t, quizL.Quizzes, 1)
	require.Equal(t, quizUUID, quizL.Quizzes[0].UUID)
	require.Equal(t, score, quizL.Quizzes[0].Score)
	require.Equal(t, "x", quizL.Quizzes[0].CorrectAnswer)
}

func TestQuizHandler_Delete(t *testing.T) {
	t.Parallel()

	token, err := jwtutils.GenerateToken(uuid.Nil, true)
	require.NoError(t, err)

	i := testutils.NewTestInjector(t,
		repositories.RepositoryPackage,
	)
	do.ProvideValue(i, slog.Default())
	do.Provide(i, testrunner.NewTestRunner)

	var quizUUID uuid.UUID
	score := 1
	db := do.MustInvoke[*bun.DB](i)
	err = db.NewInsert().Model(&models.Quiz{
		Path:          "/home/ETS/programming/Fathom/application/internal/testutils/placebo.md",
		Checksum:      []byte{},
		Score:         score,
		CorrectAnswer: "x",
	}).Returning("uuid").
		Scan(t.Context(), &quizUUID)
	require.NoError(t, err)

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()

		i = i.Scope("valid")

		do.Provide(i, handler.NewTestService)
		router, err := server.ChiServer(i)
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/v1/quiz/%v", quizUUID),
			nil,
		)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		bodyString := rec.Body.String()
		require.Equal(t, http.StatusNoContent, rec.Code, bodyString)
	})

	t.Run("Invalid", func(t *testing.T) {
		t.Parallel()

		i = i.Scope("invalid")

		do.Provide(i, handler.NewTestService)
		router, err := server.ChiServer(i)
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/v1/quiz/%v", uuid.Max),
			nil,
		)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		bodyString := rec.Body.String()
		require.Equal(t, http.StatusNotFound, rec.Code, bodyString)
	})
}

func TestQuizHandler_Put(t *testing.T) {
	t.Parallel()

	token, err := jwtutils.GenerateToken(uuid.Nil, true)
	require.NoError(t, err)

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()

		i := testutils.NewTestInjector(t,
			repositories.RepositoryPackage,
		)
		do.ProvideValue(i, slog.Default())
		do.Provide(i, testrunner.NewTestRunner)

		path := "/home/ETS/programming/Fathom/application/internal/testutils/put.md"
		var quizUUID uuid.UUID
		score := 1
		db := do.MustInvoke[*bun.DB](i)
		err = db.NewInsert().Model(&models.Quiz{
			Path:          path,
			Checksum:      []byte{},
			Score:         score,
			CorrectAnswer: "x",
		}).Returning("uuid").
			Scan(t.Context(), &quizUUID)
		require.NoError(t, err)

		do.Provide(i, handler.NewTestService)
		router, err := server.ChiServer(i)
		require.NoError(t, err)

		check := rand.Text()
		body :=
			fmt.Sprintf(`# quiz!
		there is a body! %v
		[yeah!]`, check)
		reqJSON, _ := json.Marshal(contracts.PutQuizRequest{
			Body: body,
			Meta: quiz.Frontmatter{
				Kind:  quiz.Input,
				Score: 1,
			},
		})
		req := httptest.NewRequest(
			http.MethodPut,
			fmt.Sprintf("/api/v1/quiz/%v", quizUUID),
			bytes.NewReader(reqJSON),
		)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		bodyString := rec.Body.String()
		require.Equal(t, http.StatusNoContent, rec.Code, bodyString)

		require.FileExists(t, path)
		b, err := os.ReadFile(path)
		require.NoError(t, err)

		require.Contains(t, string(b), body, string(b))
	})

	t.Run("Orphans", func(t *testing.T) {
		meta := quiz.Frontmatter{
			Kind:  quiz.Input,
			Score: 1,
		}
		testCases := []struct {
			desc string
			name string
			body string
			meta quiz.Frontmatter
		}{
			{
				desc: "No name",
				name: "",
				body: rand.Text(),
				meta: meta,
			},
			{
				desc: "No body",
				name: rand.Text(),
				body: "",
				meta: meta,
			},
			{
				desc: "No meta",
				name: rand.Text(),
				body: rand.Text(),
				meta: quiz.Frontmatter{},
			},
			{
				desc: "Nothing",
				name: "",
				body: "",
				meta: quiz.Frontmatter{},
			},
		}
		for _, tC := range testCases {
			t.Run(tC.desc, func(t *testing.T) {
				t.Parallel()

				i := testutils.NewTestInjector(t,
					repositories.RepositoryPackage,
				)
				do.ProvideValue(i, slog.Default())
				do.Provide(i, testrunner.NewTestRunner)

				do.Provide(i, handler.NewTestService)
				router, err := server.ChiServer(i)
				require.NoError(t, err)

				var quizUUID uuid.UUID
				score := 1
				db := do.MustInvoke[*bun.DB](i)
				err = db.NewInsert().Model(&models.Quiz{
					Path:          "/home/ETS/programming/Fathom/application/internal/testutils/placebo.md",
					Checksum:      []byte{},
					Score:         score,
					CorrectAnswer: "x",
				}).Returning("uuid").
					Scan(t.Context(), &quizUUID)
				require.NoError(t, err)

				reqJSON, _ := json.Marshal(contracts.PutQuizRequest{
					Body: tC.body,
					Meta: tC.meta,
				})
				req := httptest.NewRequest(
					http.MethodPut,
					fmt.Sprintf("/api/v1/quiz/%v", quizUUID),
					bytes.NewReader(reqJSON),
				)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

				rec := httptest.NewRecorder()

				router.ServeHTTP(rec, req)

				bodyString := rec.Body.String()
				require.Equal(t, http.StatusBadRequest, rec.Code, bodyString)
			})
		}
	})
}

func TestQuizHandler_Patch(t *testing.T) {
	t.Parallel()

	token, err := jwtutils.GenerateToken(uuid.Nil, true)
	require.NoError(t, err)

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()

		i := testutils.NewTestInjector(t,
			repositories.RepositoryPackage,
		)
		do.ProvideValue(i, slog.Default())
		do.Provide(i, testrunner.NewTestRunner)

		file, err := os.CreateTemp(".", "*_tmp.md")
		require.NoError(t, err)
		defer file.Close()
		// defer os.Remove(file.Name())

		template, err := os.ReadFile("/home/ETS/programming/Fathom/application/internal/testutils/placebo.md")
		require.NoError(t, err)

		_, err = file.Write(template)

		path, err := filepath.Abs(file.Name())
		require.NoError(t, err)

		var quizUUID uuid.UUID
		score := 1
		db := do.MustInvoke[*bun.DB](i)
		err = db.NewInsert().Model(&models.Quiz{
			Path:          path,
			Checksum:      []byte{},
			Score:         score,
			CorrectAnswer: "x",
		}).Returning("uuid").
			Scan(t.Context(), &quizUUID)
		require.NoError(t, err)

		do.Provide(i, handler.NewTestService)
		router, err := server.ChiServer(i)
		require.NoError(t, err)

		newName := rand.Text()
		newScore := mrand.IntN(25) + 1

		reqJSON, _ := json.Marshal(contracts.PatchQuizRequest{
			Name:  &newName,
			Score: &newScore,
		})
		req := httptest.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("/api/v1/quiz/%v", quizUUID),
			bytes.NewReader(reqJSON),
		)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNoContent, rec.Code)

		newPath := config.TurnToAbs(newName)
		require.FileExists(t, newPath)
		defer os.Remove(newPath)

		var scoreCheck int
		err = db.NewSelect().Model((*models.Quiz)(nil)).
			Where("uuid = ?", quizUUID).Column("score").
			Scan(t.Context(), &scoreCheck)

		require.Equal(t, newScore, scoreCheck)
	})

	t.Run("Orphans", func(t *testing.T) {

		testCases := []struct {
			desc  string
			name  string
			score int
		}{
			{
				desc:  "No name",
				name:  "",
				score: 1,
			},
			{
				desc:  "No score",
				name:  rand.Text(),
				score: 0,
			},
			{
				desc:  "Nothing",
				name:  "",
				score: 0,
			},
		}
		for _, tC := range testCases {
			t.Run(tC.desc, func(t *testing.T) {
				t.Parallel()

				i := testutils.NewTestInjector(t,
					repositories.RepositoryPackage,
				)
				do.ProvideValue(i, slog.Default())
				do.Provide(i, testrunner.NewTestRunner)

				do.Provide(i, handler.NewTestService)
				router, err := server.ChiServer(i)
				require.NoError(t, err)

				file, err := os.CreateTemp(".", "*_tmp.md")
				require.NoError(t, err)
				defer file.Close()
				defer os.Remove(file.Name())

				path, err := filepath.Abs(file.Name())
				require.NoError(t, err)

				var quizUUID uuid.UUID
				score := 1
				db := do.MustInvoke[*bun.DB](i)
				err = db.NewInsert().Model(&models.Quiz{
					Path:          path,
					Checksum:      []byte{},
					Score:         score,
					CorrectAnswer: "x",
				}).Returning("uuid").
					Scan(t.Context(), &quizUUID)
				require.NoError(t, err)

				nameJ := new(string)
				if tC.name != "" {
					nameJ = &tC.name
				}

				scoreJ := new(int)
				if tC.score != 0 {
					scoreJ = &tC.score
				}

				reqJSON, _ := json.Marshal(contracts.PatchQuizRequest{
					Name:  nameJ,
					Score: scoreJ,
				})
				req := httptest.NewRequest(
					http.MethodPatch,
					fmt.Sprintf("/api/v1/quiz/%v", quizUUID),
					bytes.NewReader(reqJSON),
				)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

				rec := httptest.NewRecorder()

				router.ServeHTTP(rec, req)

				require.Equal(t, http.StatusNoContent, rec.Code)

				if nameJ != nil {
					newPath := config.TurnToAbs(*nameJ)
					require.FileExists(t, newPath)
					os.Remove(newPath)
				}
			})
		}
	})

	t.Run("Zeros", func(t *testing.T) {
		t.Run("Name", func(t *testing.T) {
			t.Parallel()

			i := testutils.NewTestInjector(t,
				repositories.RepositoryPackage,
			)
			do.ProvideValue(i, slog.Default())
			do.Provide(i, testrunner.NewTestRunner)

			file, err := os.CreateTemp(".", "*_tmp.md")
			require.NoError(t, err)
			defer file.Close()
			defer os.Remove(file.Name())

			template, err := os.ReadFile("/home/ETS/programming/Fathom/application/internal/testutils/placebo.md")
			require.NoError(t, err)

			_, err = file.Write(template)

			path, err := filepath.Abs(file.Name())
			require.NoError(t, err)

			var quizUUID uuid.UUID
			score := 1
			db := do.MustInvoke[*bun.DB](i)
			err = db.NewInsert().Model(&models.Quiz{
				Path:          path,
				Checksum:      []byte{},
				Score:         score,
				CorrectAnswer: "x",
			}).Returning("uuid").
				Scan(t.Context(), &quizUUID)
			require.NoError(t, err)

			do.Provide(i, handler.NewTestService)
			router, err := server.ChiServer(i)
			require.NoError(t, err)

			newName := ""
			newScore := mrand.IntN(25) + 1

			reqJSON, _ := json.Marshal(contracts.PatchQuizRequest{
				Name:  &newName,
				Score: &newScore,
			})
			req := httptest.NewRequest(
				http.MethodPatch,
				fmt.Sprintf("/api/v1/quiz/%v", quizUUID),
				bytes.NewReader(reqJSON),
			)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, http.StatusNoContent, rec.Code)
		})

		t.Run("Score", func(t *testing.T) {
			t.Parallel()

			i := testutils.NewTestInjector(t,
				repositories.RepositoryPackage,
			)
			do.ProvideValue(i, slog.Default())
			do.Provide(i, testrunner.NewTestRunner)

			file, err := os.CreateTemp(".", "*_tmp.md")
			require.NoError(t, err)
			defer file.Close()
			defer os.Remove(file.Name())

			template, err := os.ReadFile("/home/ETS/programming/Fathom/application/internal/testutils/placebo.md")
			require.NoError(t, err)

			_, err = file.Write(template)

			path, err := filepath.Abs(file.Name())
			require.NoError(t, err)

			var quizUUID uuid.UUID
			score := 1
			db := do.MustInvoke[*bun.DB](i)
			err = db.NewInsert().Model(&models.Quiz{
				Path:          path,
				Checksum:      []byte{},
				Score:         score,
				CorrectAnswer: "x",
			}).Returning("uuid").
				Scan(t.Context(), &quizUUID)
			require.NoError(t, err)

			do.Provide(i, handler.NewTestService)
			router, err := server.ChiServer(i)
			require.NoError(t, err)

			newName := "normalName"
			newScore := 0

			reqJSON, _ := json.Marshal(contracts.PatchQuizRequest{
				Name:  &newName,
				Score: &newScore,
			})
			req := httptest.NewRequest(
				http.MethodPatch,
				fmt.Sprintf("/api/v1/quiz/%v", quizUUID),
				bytes.NewReader(reqJSON),
			)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, http.StatusNoContent, rec.Code)
		})
	})

	t.Run("Invalid score", func(t *testing.T) {
		t.Parallel()

		i := testutils.NewTestInjector(t,
			repositories.RepositoryPackage,
		)
		do.ProvideValue(i, slog.Default())
		do.Provide(i, testrunner.NewTestRunner)

		file, err := os.CreateTemp(".", "*_tmp.md")
		require.NoError(t, err)
		defer file.Close()
		defer os.Remove(file.Name())

		template, err := os.ReadFile("/home/ETS/programming/Fathom/application/internal/testutils/placebo.md")
		require.NoError(t, err)

		_, err = file.Write(template)

		path, err := filepath.Abs(file.Name())
		require.NoError(t, err)

		var quizUUID uuid.UUID
		score := 1
		db := do.MustInvoke[*bun.DB](i)
		err = db.NewInsert().Model(&models.Quiz{
			Path:          path,
			Checksum:      []byte{},
			Score:         score,
			CorrectAnswer: "x",
		}).Returning("uuid").
			Scan(t.Context(), &quizUUID)
		require.NoError(t, err)

		do.Provide(i, handler.NewTestService)
		router, err := server.ChiServer(i)
		require.NoError(t, err)

		newName := "normalName"
		newScore := -5423

		reqJSON, _ := json.Marshal(contracts.PatchQuizRequest{
			Name:  &newName,
			Score: &newScore,
		})
		req := httptest.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("/api/v1/quiz/%v", quizUUID),
			bytes.NewReader(reqJSON),
		)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)

	})
}
