package importutils

import (
	"fmt"
	"net/http"

	"github.com/egot3/fathom/internal/carefulness"
	exportutils "github.com/egot3/fathom/internal/exportUtils"
	"github.com/egot3/fathom/internal/models"
)

// thingabob which assigns quiz to test
func reconcile(quizzes []StagedQuiz, m *exportutils.YamlTest) ([]StagedQuiz, StagedTest, carefulness.JSONErrorable) {
	if m == nil {
		return quizzes, StagedTest{}, nil
	}

	byPath := make(map[string]int, len(quizzes))
	for i, q := range quizzes {
		byPath[q.RelPath] = i
	}

	test := StagedTest{Test: models.Test{UUID: m.UUID, Name: m.Name}}
	for _, mq := range m.Quizzes {
		i, inArchive := byPath[cleanName(mq.Path)]
		if !inArchive {
			test.QuizUUIDs = append(test.QuizUUIDs, mq.UUID)
			continue
		}
		if got := quizzes[i].Quiz.Checksum.String(); got != mq.Checksum {
			return nil, StagedTest{}, carefulness.JSONError{
				Err:    fmt.Sprintf("checksum mismatch for %s", mq.Path),
				Status: http.StatusUnprocessableEntity,
			}
		}
		quizzes[i].Quiz.UUID = mq.UUID
		test.QuizUUIDs = append(test.QuizUUIDs, mq.UUID)
	}
	return quizzes, test, nil
}
