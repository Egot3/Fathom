package testutils

import (
	"embed"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed placebo.md
var placebo embed.FS

func TestQuiz(t *testing.T) *os.File {
	t.Helper()

	f, err := os.CreateTemp("", "*.md")
	require.NoError(t, err)

	pb, err := placebo.ReadFile("placebo.md")
	require.NoError(t, err)

	_, err = f.Write(pb)
	require.NoError(t, err)

	return f
}
