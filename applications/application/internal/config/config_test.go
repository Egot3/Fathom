package config_test

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/egot3/fathom/internal/config"
	"github.com/stretchr/testify/require"
)

func TestParseLogSinks(t *testing.T) {
	require.Equal(t, []string{"slog"}, config.ParseLogSinks(""))
	require.Equal(t, []string{"slog"}, config.ParseLogSinks("   "))
	require.Equal(t, []string{"slog", "charmLog"}, config.ParseLogSinks("slog,charmLog"))
}

func TestTurnToAbs(t *testing.T) {
	c := config.Config{QuizPath: "/data/quizzes"}
	p, err := c.TurnToAbs("some/quiz")
	require.NoError(t, err)
	require.Equal(t, filepath.Join(filepath.VolumeName(os.Args[0]), "/data/quizzes/some/quiz.md"), p)
}

func TestGetEnv(t *testing.T) {
	t.Run("Has one", func(t *testing.T) {
		key := rand.Text()
		value := rand.Text()
		os.Setenv(key, value)

		v1 := config.GetEnv(key, "", nil)
		require.Equal(t, value, v1)
	})
	t.Run("Fit constraints", func(t *testing.T) {
		key := rand.Text()
		value := rand.Text()
		os.Setenv(key, value)

		v1 := config.GetEnv(key, "", []string{value, rand.Text()})
		require.Equal(t, value, v1)
	})
	t.Run("Default", func(t *testing.T) {
		key := rand.Text()
		value := rand.Text()

		v1 := config.GetEnv(key, value, nil)
		require.Equal(t, value, v1)
	})
	t.Run("Doesn't fit constraints", func(t *testing.T) {
		key := rand.Text()
		value := rand.Text()
		os.Setenv(key, value)

		constr := rand.Text()
		v1 := config.GetEnv(key, constr, []string{constr})
		require.Equal(t, constr, v1)
	})

}
