package importutils

import (
	"fmt"
	"net/http"

	"github.com/egot3/fathom/internal/carefulness"
	exportutils "github.com/egot3/fathom/internal/exportUtils"
	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
)

// thingabob which assigns quiz to test
func reconcile(quizzes []StagedQuiz, m *exportutils.Manifest) ([]StagedQuiz, StagedTest, carefulness.JSONErrorable) {
	if m == nil {
		return quizzes, StagedTest{}, nil
	}

	byPath := make(map[string]int, len(quizzes))
	for i, q := range quizzes {
		byPath[q.RelPath] = i
	}

	var test StagedTest
	addToStaged := func(UUID uuid.UUID) {}
	if m.Kind == exportutils.Test {
		test = StagedTest{Test: models.Test{UUID: m.UUID, Name: m.Name}}
		addToStaged = func(UUID uuid.UUID) {
			test.QuizUUIDs = append(test.QuizUUIDs, UUID)
		}
	}

	for _, mq := range m.Quizzes {
		i, inArchive := byPath[cleanName(mq.Path)]
		if !inArchive {
			addToStaged(mq.UUID)
			continue
		}
		if got := quizzes[i].Quiz.Checksum.String(); got != mq.Checksum {
			return nil, StagedTest{}, carefulness.JSONError{
				Err:    fmt.Sprintf("checksum mismatch for %s", mq.Path),
				Status: http.StatusUnprocessableEntity,
			}
		}
		quizzes[i].Quiz.UUID = mq.UUID
		addToStaged(mq.UUID)
	}
	return quizzes, test, nil
}
