package handler_test

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/egot3/fathom/internal/contracts"
	"github.com/egot3/fathom/internal/database/repositories"
	"github.com/egot3/fathom/internal/handler"
	"github.com/egot3/fathom/internal/models"
	testrunner "github.com/egot3/fathom/internal/testRunner"
	"github.com/egot3/fathom/internal/testutils"
	"github.com/egot3/fathom/server"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

func TestUserHandler_Register(t *testing.T) {
	t.Parallel()

	t.Run("General test cases", func(t *testing.T) {
		t.Parallel()

		generalTestCases := []struct {
			desc       string
			wantStatus int
			nickname   string
			password   string
		}{
			{
				desc:       "Valid registration",
				password:   "1253Eq13",
				wantStatus: http.StatusCreated,
				nickname:   rand.Text(),
			},
			{
				desc:       "No nickname",
				nickname:   "",
				password:   "12254Eq13",
				wantStatus: http.StatusBadRequest,
			},
			{
				desc:       "No password",
				nickname:   rand.Text(),
				password:   "",
				wantStatus: http.StatusBadRequest,
			},
		}
		for _, gtt := range generalTestCases {
			t.Run(gtt.desc, func(t *testing.T) {
				t.Parallel()

				i := testutils.NewTestInjector(t,
					repositories.RepositoryPackage,
				)
				do.ProvideValue(i, slog.Default())
				do.Provide(i, testrunner.NewTestRunner)

				do.Provide(i, handler.NewTestService)
				router, err := server.ChiServer(i)
				require.NoError(t, err)

				req := httptest.NewRequest(
					http.MethodPost,
					"/api/v1/user/register",
					strings.NewReader(fmt.Sprintf(`{"nickname": %q, "password": %q}`, gtt.nickname, gtt.password)),
				)
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				router.ServeHTTP(rec, req)

				bodyString := rec.Body.String()
				require.Equal(t, gtt.wantStatus, rec.Code)
				if gtt.wantStatus < 300 {
					require.Contains(t, bodyString, gtt.nickname, "Bad nickname")
				}

			})
		}
	})

	t.Run("Password test cases", func(t *testing.T) {
		passwordTestCases := []struct {
			desc        string
			password    string
			errContains []string
		}{
			{
				desc:        "Whitespace",
				password:    "1254Eq1 3",
				errContains: []string{"mustn't", "whitespace"},
			},
			{
				desc:        "\\n",
				password:    "1254Eq1\n3",
				errContains: []string{"mustn't", "whitespace"},
			},
			{
				desc:        "\\r",
				password:    "1254Eq1\r3",
				errContains: []string{"mustn't", "whitespace"},
			},
			{
				desc:        "\\t",
				password:    "1254Eq1\t3",
				errContains: []string{"mustn't", "whitespace"},
			},
			{
				desc:        "Too little lowercase",
				password:    "1254EQ13",
				errContains: []string{"must", "at least 1", "lowercase"},
			},
			{
				desc:        "Too little uppercase",
				password:    "1254eq13",
				errContains: []string{"must", "at least 1", "uppercase"},
			},
			{
				desc:        "Too short",
				password:    "1254e",
				errContains: []string{"must", "at least 8"},
			},
			{
				desc:        "Too long",
				password:    "1254eq1dsadasdaEwfsadasdadasdcZXCsxsadxZCZXXASD31251254eq1dsadasdaEwfsadasdadasdcZXCsxsadxZCZXXASD31254eq1dsadasdaEwfsadasdadasdcZXCsxsadxZCZXXASD34eq1dsadasdaEwfsadasdadasdcZXCsxsadxZCZXXASD3",
				errContains: []string{"mustn't", "be longer", "60"},
			},
			{
				desc:        "Too little numbers",
				password:    "12ErQeqwsdEq",
				errContains: []string{"must", "at least", "6", "numbers"},
			},
		}
		for _, ptt := range passwordTestCases {
			t.Run(ptt.desc, func(t *testing.T) {
				t.Parallel()

				i := testutils.NewTestInjector(t,
					repositories.RepositoryPackage,
				)
				do.ProvideValue(i, slog.Default())

				do.Provide(i, testrunner.NewTestRunner)

				do.Provide(i, handler.NewTestService)
				router, err := server.ChiServer(i)
				require.NoError(t, err)

				req := httptest.NewRequest(
					http.MethodPost,
					"/api/v1/user/register",
					strings.NewReader(fmt.Sprintf(`{"nickname": %q, "password": %q}`, rand.Text(), ptt.password)),
				)
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				router.ServeHTTP(rec, req)

				require.Equal(t, http.StatusBadRequest, rec.Code)
				bodyString := rec.Body.String()
				require.Contains(t, bodyString, "error")
				for _, cont := range ptt.errContains {
					require.Contains(t, bodyString, cont)
				}
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

		do.Provide(i, handler.NewTestService)
		router, err := server.ChiServer(i)
		require.NoError(t, err)

		nickname := rand.Text()
		req := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/user/register",
			strings.NewReader(fmt.Sprintf(`{"nickname": %q, "password": %q}`, nickname, "123123EqdsaE")),
		)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		bodyString := rec.Body.String()
		require.Equal(t, http.StatusCreated, rec.Code)
		require.Contains(t, bodyString, nickname, "Bad nickname")

		conflictingReq := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/user/register",
			strings.NewReader(fmt.Sprintf(`{"nickname": %q, "password": %q}`, nickname, "123123EqdsasE")),
		)
		rec = httptest.NewRecorder()

		router.ServeHTTP(rec, conflictingReq)
		bodyString = rec.Body.String()
		require.Equal(t, http.StatusConflict, rec.Code)
	})
}

func TestUserHandler_Login(t *testing.T) {
	t.Parallel()

	i := testutils.NewTestInjector(t,
		repositories.RepositoryPackage,
	)
	do.ProvideValue(i, slog.Default())
	do.Provide(i, testrunner.NewTestRunner)

	db := do.MustInvoke[*bun.DB](i)

	pswd := rand.Text()
	pswdhash, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := models.User{Nickname: rand.Text(), PasswordHash: pswdhash}
	err = db.NewInsert().Model(&user).Scan(t.Context())
	require.NoError(t, err)

	testCases := []struct {
		desc       string
		nickname   string
		password   string
		wantStatus int
	}{
		{
			desc:       "Valid login",
			nickname:   user.Nickname,
			password:   pswd,
			wantStatus: http.StatusOK,
		},
		{
			desc:       "No nickname",
			nickname:   "",
			password:   pswd,
			wantStatus: http.StatusBadRequest,
		},
		{
			desc:       "No password",
			nickname:   user.Nickname,
			password:   "",
			wantStatus: http.StatusBadRequest,
		},
		{
			desc:       "Wrong nickname",
			nickname:   "wrong!",
			password:   pswd,
			wantStatus: http.StatusNotFound,
		},
		{
			desc:       "Wrong password",
			nickname:   user.Nickname,
			password:   "wrong password",
			wantStatus: http.StatusNotFound,
		},
	}
	for _, ltt := range testCases {
		t.Run(ltt.desc, func(t *testing.T) {
			t.Parallel()

			i := i.Scope(ltt.desc)

			do.Provide(i, handler.NewTestService)
			router, err := server.ChiServer(i)
			require.NoError(t, err)

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/v1/user/login",
				strings.NewReader(fmt.Sprintf(`{"nickname": %q, "password": %q}`, ltt.nickname, ltt.password)),
			)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, ltt.wantStatus, rec.Code)

		})
	}

}

func TestUserHandler_Get(t *testing.T) {
	t.Parallel()

	i := testutils.NewTestInjector(t,
		repositories.RepositoryPackage,
	)
	do.ProvideValue(i, slog.Default())
	do.Provide(i, testrunner.NewTestRunner)

	db := do.MustInvoke[*bun.DB](i)

	pswd := rand.Text()
	pswdhash, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := models.User{Nickname: rand.Text(), PasswordHash: pswdhash}
	err = db.NewInsert().Model(&user).Scan(t.Context())
	require.NoError(t, err)

	t.Run("Found", func(t *testing.T) {
		t.Parallel()

		i := i.Scope(t.Name())

		do.Provide(i, handler.NewTestService)
		router, err := server.ChiServer(i)
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/user/%v", user.UUID),
			strings.NewReader(""),
		)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, 200, rec.Code)

		var contract contracts.GetUserResponse
		err = json.Unmarshal(rec.Body.Bytes(), &contract)
		require.NoError(t, err)

		retUser := contract.User
		require.Equal(t, user.Nickname, retUser.Nickname)
	})
	t.Run("Not found", func(t *testing.T) {
		t.Parallel()

		i := i.Scope(t.Name())

		do.Provide(i, handler.NewTestService)
		router, err := server.ChiServer(i)
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/user/%v", uuid.Nil),
			nil,
		)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, 404, rec.Code)
	})
}

func TestUserHandler_List(t *testing.T) {
	t.Parallel()

	i := testutils.NewTestInjector(t,
		repositories.RepositoryPackage,
	)
	do.ProvideValue(i, slog.Default())
	do.Provide(i, testrunner.NewTestRunner)

	db := do.MustInvoke[*bun.DB](i)

	users := make([]models.User, 5)
	for i := range 5 {
		pswd := rand.Text()
		pswdhash, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.DefaultCost)
		require.NoError(t, err)

		users[i] = models.User{PasswordHash: pswdhash, Nickname: rand.Text()}
	}
	err := db.NewInsert().Model(&users).Scan(t.Context())
	require.NoError(t, err)

	t.Run("Listing users", func(t *testing.T) {
		t.Parallel()

		do.Provide(i, handler.NewTestService)
		router, err := server.ChiServer(i)
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/user?size=%v&page=%v", 5, 0),
			nil,
		)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, 200, rec.Code)

		var contract contracts.ListUsersResponse
		err = json.Unmarshal(rec.Body.Bytes(), &contract)
		require.NoError(t, err)

		retUser := contract.Users
		for _, user := range users {
			require.True(t, lo.SomeBy(retUser, func(retUser models.User) bool {
				return retUser.UUID == user.UUID
			}))
		}
	})
}
