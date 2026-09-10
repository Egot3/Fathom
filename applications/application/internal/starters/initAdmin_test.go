package starters_test

import (
	"crypto/rand"
	"testing"

	"github.com/egot3/fathom/internal/config"
	"github.com/egot3/fathom/internal/models"
	"github.com/egot3/fathom/internal/starters"
	"github.com/egot3/fathom/internal/testutils"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

func TestInitAdmin(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		i := testutils.NewTestInjector(t)
		do.OverrideValue(i, &config.Config{LogLevel: "debug"})

		db := do.MustInvoke[*bun.DB](i)

		starters.InitAdmin(i)

		count, err := db.NewSelect().Model((*models.User)(nil)).Count(t.Context())
		require.NoError(t, err)

		require.Equal(t, 0, count)
	})

	t.Run("Not empty", func(t *testing.T) {
		name := rand.Text()
		password := rand.Text()

		i := testutils.NewTestInjector(t)
		do.OverrideValue(i, &config.Config{
			InitAdminUsername: name,
			InitAdminPassword: password,
			LogLevel:          "debug",
		})

		db := do.MustInvoke[*bun.DB](i)

		starters.InitAdmin(i)

		var admin models.User
		count, err := db.NewSelect().Model(&admin).ScanAndCount(t.Context())
		require.NoError(t, err)

		require.Equal(t, 1, count)
		require.Equal(t, name, admin.Nickname)
		require.True(t, admin.IsTeacher)
		require.NoError(t, bcrypt.CompareHashAndPassword(admin.PasswordHash, []byte(password)))
	})

	t.Run("Exists", func(t *testing.T) {
		t.Run("Override", func(t *testing.T) {
			t.Run("Promote", func(t *testing.T) {
				name := rand.Text()
				password := rand.Text()

				i := testutils.NewTestInjector(t)
				do.OverrideValue(i, &config.Config{
					InitAdminUsername: name,
					InitAdminPassword: rand.Text(),
					LogLevel:          "debug",
				})

				db := do.MustInvoke[*bun.DB](i)

				pswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				require.NoError(t, err)

				_, err = db.NewInsert().Model(&models.User{Nickname: name, PasswordHash: pswd}).Exec(t.Context())
				require.NoError(t, err)

				starters.InitAdmin(i)

				var admin models.User
				count, err := db.NewSelect().Model(&admin).ScanAndCount(t.Context())
				require.NoError(t, err)

				require.Equal(t, 1, count)
				require.Equal(t, name, admin.Nickname)
				require.True(t, admin.IsTeacher)
				require.NotEqual(t, pswd, admin.PasswordHash)
			})
			t.Run("Repassword", func(t *testing.T) {
				name := rand.Text()
				password := rand.Text()

				i := testutils.NewTestInjector(t)
				do.OverrideValue(i, &config.Config{
					InitAdminUsername: name,
					InitAdminPassword: rand.Text(),
					LogLevel:          "debug",
				})

				db := do.MustInvoke[*bun.DB](i)

				pswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				require.NoError(t, err)

				_, err = db.NewInsert().Model(&models.User{Nickname: name, PasswordHash: pswd, IsTeacher: true}).Exec(t.Context())
				require.NoError(t, err)

				starters.InitAdmin(i)

				var admin models.User
				count, err := db.NewSelect().Model(&admin).ScanAndCount(t.Context())
				require.NoError(t, err)

				require.Equal(t, 1, count)
				require.Equal(t, name, admin.Nickname)
				require.True(t, admin.IsTeacher)
				require.NotEqual(t, pswd, admin.PasswordHash)
			})
		})

		t.Run("Immutable", func(t *testing.T) {
			name := rand.Text()
			password := rand.Text()

			i := testutils.NewTestInjector(t)
			do.OverrideValue(i, &config.Config{
				LogLevel: "debug",
			})

			db := do.MustInvoke[*bun.DB](i)

			pswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			require.NoError(t, err)

			_, err = db.NewInsert().Model(&models.User{Nickname: name, PasswordHash: pswd, IsTeacher: true}).Exec(t.Context())
			require.NoError(t, err)

			starters.InitAdmin(i)

			var admin models.User
			count, err := db.NewSelect().Model(&admin).ScanAndCount(t.Context())
			require.NoError(t, err)

			require.Equal(t, 1, count)
			require.Equal(t, name, admin.Nickname)
			require.True(t, admin.IsTeacher)
			require.Equal(t, pswd, admin.PasswordHash)
		})
	})
}
