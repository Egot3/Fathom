package quizparser_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	quizparser "github.com/egot3/fathom/internal/quizParser"
	"github.com/stretchr/testify/require"
)

func Test_parsers(t *testing.T) {
	root := "../../quizzes/test"
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			t.Logf("Directory: %s\n", path)
		} else {
			t.Run(path, func(t *testing.T) {
				raw, err := os.ReadFile(path)
				require.NoError(t, err)
				q, err := quizparser.ParseQuizByBytes(raw)
				if strings.Contains(path, "valid") {
					require.NoError(t, err)
					t.Logf("%+v", q)
				} else {
					t.Logf("path:%v\n%v", path, err)
					require.Error(t, err)
				}
			})

		}

		return nil
	})
	require.NoError(t, err)
}
