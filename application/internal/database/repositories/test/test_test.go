package test_test //finnaly something remotely funny

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"io"
	"log"
	mrand "math/rand/v2"
	"testing"

	"github.com/egot3/fathom/internal/database/repositories/test"
	"github.com/egot3/fathom/internal/models"
	"github.com/egot3/fathom/internal/testutils"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func NewInjectorWithTestRepo(t testing.TB) do.Injector {
	t.Helper()

	i := testutils.NewTestInjector(t)

	do.Provide(i, test.NewTestRepository)

	return i
}

func RegisterModels(db *bun.DB) {
	db.RegisterModel((*models.TestsQuizzies)(nil))
	db.RegisterModel((*models.GroupsUsers)(nil))
}

func TestTest_Creation(t *testing.T) {
	t.Parallel()

	i := NewInjectorWithTestRepo(t)
	r := do.MustInvoke[test.TestRepository](i)
	db := do.MustInvoke[*bun.DB](i)

	RegisterModels(db)

	t.Run("New test", func(t *testing.T) {
		t.Parallel()

		name := rand.Text()
		err := r.CreateTest(t.Context(), name)
		require.NoError(t, err)

		var test models.Test
		err = db.NewSelect().Model(&test).Where("name = ?", name).Scan(t.Context())
		require.NoError(t, err)
		require.Equal(t, name, test.Name)
	})

	t.Run("Existing test", func(t *testing.T) {
		t.Parallel()

		name := rand.Text()
		err := r.CreateTest(t.Context(), name)
		require.NoError(t, err)

		var test models.Test
		err = db.NewSelect().Model(&test).Where("name = ?", name).Scan(t.Context())
		require.NoError(t, err)
		require.Equal(t, name, test.Name)

		err = r.CreateTest(t.Context(), name)
		require.Error(t, err)
	})
}

func TestTest_Deletion(t *testing.T) {
	t.Parallel()

	i := NewInjectorWithTestRepo(t)
	r := do.MustInvoke[test.TestRepository](i)
	db := do.MustInvoke[*bun.DB](i)

	RegisterModels(db)

	test := models.Test{Name: rand.Text()}
	err := db.NewInsert().Model(&test).Returning("*").Scan(t.Context())
	require.NoError(t, err)

	t.Run("Existing test", func(t *testing.T) {
		t.Parallel()

		err := r.DeleteTest(t.Context(), test.UUID)
		require.NoError(t, err)

		err = db.NewSelect().Model(&test).WherePK().Scan(t.Context())
		require.Error(t, err, test)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("Non-existing test", func(t *testing.T) {
		t.Parallel()

		err := r.DeleteTest(t.Context(), uuid.Nil)
		require.Error(t, err, test)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestTest_Update(t *testing.T) {
	t.Parallel()

	i := NewInjectorWithTestRepo(t)
	r := do.MustInvoke[test.TestRepository](i)
	db := do.MustInvoke[*bun.DB](i)

	RegisterModels(db)

	test := models.Test{Name: rand.Text()}
	err := db.NewInsert().Model(&test).Returning("*").Scan(t.Context())
	require.NoError(t, err)

	t.Run("Found test", func(t *testing.T) {
		t.Parallel()

		err := r.UpdateTest(t.Context(), test.UUID, rand.Text())
		require.NoError(t, err)

		testR := models.Test{UUID: test.UUID}
		err = db.NewSelect().Model(&testR).WherePK().Scan(t.Context())
		require.NoError(t, err)

		require.NotEqual(t, test.Name, testR.Name)
	})
	t.Run("Not found test", func(t *testing.T) {
		t.Parallel()

		err := r.UpdateTest(t.Context(), uuid.Nil, rand.Text())
		require.Error(t, err)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestTest_Read(t *testing.T) {
	t.Parallel()

	i := NewInjectorWithTestRepo(t)
	r := do.MustInvoke[test.TestRepository](i)
	db := do.MustInvoke[*bun.DB](i)

	RegisterModels(db)

	test := models.Test{Name: rand.Text()}
	err := db.NewInsert().Model(&test).Returning("*").Scan(t.Context())
	require.NoError(t, err)

	t.Run("Found test", func(t *testing.T) {
		t.Parallel()

		testR, err := r.Test(t.Context(), test.UUID)
		require.NoError(t, err)

		require.Equal(t, test.Name, testR.Name)
		require.Equal(t, test.UUID, testR.UUID)
		require.Equal(t, test.CreatedAt, testR.CreatedAt)
		require.Equal(t, test.UpdatedAt, testR.UpdatedAt)
		require.Empty(t, testR.Quizzes)
	})

	t.Run("Not found test", func(t *testing.T) {
		t.Parallel()

		testR, err := r.Test(t.Context(), uuid.Nil)
		require.Error(t, err)
		require.ErrorIs(t, err, sql.ErrNoRows)
		require.Empty(t, testR)
	})
}

func TestTest_Bundle(t *testing.T) {
	t.Parallel()

	i := NewInjectorWithTestRepo(t)
	r := do.MustInvoke[test.TestRepository](i)
	db := do.MustInvoke[*bun.DB](i)

	RegisterModels(db)

	t.Run("Valid bundle", func(t *testing.T) {
		test := models.Test{Name: rand.Text()}
		err := db.NewInsert().Model(&test).Returning("*").Scan(t.Context())
		require.NoError(t, err)

		var quizzes []models.Quiz
		for range mrand.IntN(6) + 3 {
			quizzes = append(quizzes, models.Quiz{Path: fmt.Sprintf("/path/to/%v.md", rand.Text()), Checksum: []byte{}})
		}
		err = db.NewInsert().Model(&quizzes).Returning("*").Scan(t.Context())
		require.NoError(t, err)

		pathes := lo.Map(quizzes, func(quiz models.Quiz, _ int) string { return quiz.Path })
		err = r.BundleQuizzesToTest(t.Context(), test.UUID, pathes)
		require.NoError(t, err)

		type quizShort struct {
			Path     string `bun:"quiz_path"`
			Position int    `bun:"position"`
		}
		var quizzesR []quizShort
		err = db.NewSelect().Model((*models.TestsQuizzies)(nil)).
			Where("test_uuid = ?", test.UUID).
			Column("quiz_path", "position").
			Scan(t.Context(), &quizzesR)

		var pathesR []string
		var posR []int
		lo.ForEach(quizzesR, func(quizS quizShort, _ int) {
			pathesR = append(pathesR, quizS.Path)
			posR = append(posR, quizS.Position)

		})
		require.Condition(t, func() (success bool) {
			i := -1 //crutch
			return lo.EveryBy(pathes, func(path string) bool {
				i++
				return i == posR[lo.IndexOf(pathesR, path)]
			})
		}, pathesR, posR, pathes)
	})

	t.Run("Patially valid bundle", func(t *testing.T) {
		testM := models.Test{Name: rand.Text()}
		err := db.NewInsert().Model(&testM).Returning("*").Scan(t.Context())
		require.NoError(t, err)

		var quizzes []models.Quiz
		for range mrand.IntN(6) + 3 {
			quizzes = append(quizzes, models.Quiz{Path: fmt.Sprintf("/path/to/%v.md", rand.Text()), Checksum: []byte{}})
		}
		err = db.NewInsert().Model(&quizzes).Returning("*").Scan(t.Context())
		require.NoError(t, err)

		pathes := lo.Map(quizzes, func(quiz models.Quiz, _ int) string { return quiz.Path })
		err = r.BundleQuizzesToTest(t.Context(), testM.UUID, append(pathes, "/unknown.md"))
		require.Error(t, err)
	})

	t.Run("Invalid bundle", func(t *testing.T) {
		testM := models.Test{Name: rand.Text()}
		err := db.NewInsert().Model(&testM).Returning("*").Scan(t.Context())
		require.NoError(t, err)

		err = r.BundleQuizzesToTest(t.Context(), testM.UUID, append([]string{}, "/unknown.md"))
		require.Error(t, err)
	})
}

// --- Now making benchmarks for those, which take in slices
// as they are more prone to bad perfomance because of not optimal programm
func BenchmarkGroup_Bundle_quizzes(b *testing.B) {
	log.SetOutput(io.Discard)
	b.Run("Benchmark 5 appendants", func(b *testing.B) {
		i := NewInjectorWithTestRepo(b)

		r := do.MustInvoke[test.TestRepository](i)
		db := do.MustInvoke[*bun.DB](i)
		RegisterModels(db)

		name := rand.Text()
		testUUID := uuid.UUID{}
		err := db.NewInsert().Model(&models.Test{Name: name}).Returning("uuid").Scan(b.Context(), &testUUID)
		require.NoError(b, err)

		var pathes []string
		var quizzes []models.Quiz
		for range 5 {
			quizPath := fmt.Sprintf("/path/to/%v.md", rand.Text())

			require.NoError(b, err)
			pathes = append(pathes, quizPath)
			quizzes = append(quizzes, models.Quiz{Path: quizPath, Checksum: []byte{}})
		}

		_, err = db.NewInsert().Model(&quizzes).
			Exec(b.Context())

		b.ResetTimer()

		for b.Loop() {
			b.StopTimer()
			_, err := db.NewTruncateTable().Table("tests_quizzes").Exec(b.Context())
			require.NoError(b, err)
			b.StartTimer()

			err = r.BundleQuizzesToTest(b.Context(), testUUID, pathes)
			require.NoError(b, err)
		}
	})

	b.Run("Benchmark 50 appendants", func(b *testing.B) {
		i := NewInjectorWithTestRepo(b)

		r := do.MustInvoke[test.TestRepository](i)
		db := do.MustInvoke[*bun.DB](i)
		RegisterModels(db)

		name := rand.Text()
		testUUID := uuid.UUID{}
		err := db.NewInsert().Model(&models.Test{Name: name}).Returning("uuid").Scan(b.Context(), &testUUID)
		require.NoError(b, err)

		var pathes []string
		var quizzes []models.Quiz
		for range 50 {
			quizPath := fmt.Sprintf("/path/to/%v.md", rand.Text())

			require.NoError(b, err)
			pathes = append(pathes, quizPath)
			quizzes = append(quizzes, models.Quiz{Path: quizPath, Checksum: []byte{}})
		}

		_, err = db.NewInsert().Model(&quizzes).
			Exec(b.Context())

		b.ResetTimer()

		for b.Loop() {
			b.StopTimer()
			_, err := db.NewTruncateTable().Table("tests_quizzes").Exec(b.Context())
			require.NoError(b, err)
			b.StartTimer()

			err = r.BundleQuizzesToTest(b.Context(), testUUID, pathes)
			require.NoError(b, err)
		}
	})

	b.Run("Benchmark 500 appendants", func(b *testing.B) {
		i := NewInjectorWithTestRepo(b)

		r := do.MustInvoke[test.TestRepository](i)
		db := do.MustInvoke[*bun.DB](i)
		RegisterModels(db)

		name := rand.Text()
		testUUID := uuid.UUID{}
		err := db.NewInsert().Model(&models.Test{Name: name}).Returning("uuid").Scan(b.Context(), &testUUID)
		require.NoError(b, err)

		var pathes []string
		var quizzes []models.Quiz
		for range 500 {
			quizPath := fmt.Sprintf("/path/to/%v.md", rand.Text())

			require.NoError(b, err)
			pathes = append(pathes, quizPath)
			quizzes = append(quizzes, models.Quiz{Path: quizPath, Checksum: []byte{}})
		}

		_, err = db.NewInsert().Model(&quizzes).
			Exec(b.Context())

		b.ResetTimer()

		for b.Loop() {
			b.StopTimer()
			_, err := db.NewTruncateTable().Table("tests_quizzes").Exec(b.Context())
			require.NoError(b, err)
			b.StartTimer()

			err = r.BundleQuizzesToTest(b.Context(), testUUID, pathes)
			require.NoError(b, err)
		}
	})
}
