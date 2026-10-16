package handler_test

import (
	"bytes"
	"crypto/rand"
	"encoding/csv"
	mrand "math/rand/v2"
	"strconv"
	"testing"

	"github.com/egot3/fathom/internal/handler"
	"github.com/egot3/fathom/internal/models"
	"github.com/egot3/fathom/internal/testutils"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func TestCommitUserImport(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		i := testutils.NewTestInjector(t, handler.HandlerPackage)

		svc, err := handler.SeviceStruct(i)
		require.NoError(t, err)

		size := mrand.IntN(50) + 1
		rw := bytes.NewBuffer([]byte{})
		w := csv.NewWriter(rw)
		w.Write([]string{"name", "password"})
		for range size {
			w.Write([]string{
				rand.Text(),
				"superSAFEPASSWORD!#231312!!",
			})
		}
		w.Flush()

		r := csv.NewReader(rw)
		err = handler.CommitUserImport(t.Context(), r, svc)
		require.NoError(t, err)

		db := do.MustInvoke[*bun.DB](i)
		c, err := db.NewSelect().Model((*models.User)(nil)).Count(t.Context())
		require.NoError(t, err)
		require.EqualValues(t, size, c)
	})
	t.Run("Invalid", func(t *testing.T) {
		t.Run("No header", func(t *testing.T) {
			i := testutils.NewTestInjector(t, handler.HandlerPackage)

			svc, err := handler.SeviceStruct(i)
			require.NoError(t, err)

			size := mrand.IntN(50) + 1
			rw := bytes.NewBuffer([]byte{})
			w := csv.NewWriter(rw)
			for range size {
				w.Write([]string{
					rand.Text(),
					"superSAFEPASSWORD!#231312!!",
				})
			}
			w.Flush()

			r := csv.NewReader(rw)
			err = handler.CommitUserImport(t.Context(), r, svc)
			require.Error(t, err)

			db := do.MustInvoke[*bun.DB](i)
			c, err := db.NewSelect().Model((*models.User)(nil)).Count(t.Context())
			require.NoError(t, err)
			require.EqualValues(t, 0, c)
		})
		t.Run("Bad password", func(t *testing.T) {
			i := testutils.NewTestInjector(t, handler.HandlerPackage)

			svc, err := handler.SeviceStruct(i)
			require.NoError(t, err)

			size := mrand.IntN(50) + 1
			rw := bytes.NewBuffer([]byte{})
			w := csv.NewWriter(rw)
			w.Write([]string{"name", "password"})
			for range size {
				w.Write([]string{
					rand.Text(),
					"unsec",
				})
			}
			w.Flush()

			r := csv.NewReader(rw)
			err = handler.CommitUserImport(t.Context(), r, svc)
			require.Error(t, err)

			db := do.MustInvoke[*bun.DB](i)
			c, err := db.NewSelect().Model((*models.User)(nil)).Count(t.Context())
			require.NoError(t, err)
			require.EqualValues(t, 0, c)
		})
		t.Run("Badly formatted", func(t *testing.T) {
			i := testutils.NewTestInjector(t, handler.HandlerPackage)

			svc, err := handler.SeviceStruct(i)
			require.NoError(t, err)

			size := mrand.IntN(50) + 1
			rw := bytes.NewBuffer([]byte{})
			w := csv.NewWriter(rw)
			w.Write([]string{"name", "password", "is_teacher"})
			for range size {
				w.Write([]string{
					rand.Text(),
					"superSAFEPASSWORD!#231312!!",
					strconv.FormatBool(bool(mrand.IntN(2) == 1)),
				})
			}
			w.Flush()

			r := csv.NewReader(rw)
			err = handler.CommitUserImport(t.Context(), r, svc)
			require.Error(t, err)

			db := do.MustInvoke[*bun.DB](i)
			c, err := db.NewSelect().Model((*models.User)(nil)).Count(t.Context())
			require.NoError(t, err)
			require.EqualValues(t, 0, c)
		})
	})
}
