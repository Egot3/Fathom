package quizparser

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TestRunner struct {
	mu      sync.RWMutex
	name    string
	quizzes map[uuid.UUID]*Quiz
}

func (tr *TestRunner) Start(ctx context.Context, quizPathesWithUUIDs map[uuid.UUID]string) error {
	quizzes := make(map[uuid.UUID]*Quiz, len(quizPathesWithUUIDs))
	for uuid, path := range quizPathesWithUUIDs {
		_, exists := quizzes[uuid]
		if exists {
			continue
		}
		switch path[:2] {
		case "./":
			quiz, err := ParseQuizByPath(path)
			if err != nil {
				return err
			}
			quizzes[uuid] = quiz
		default:
			return fmt.Errorf("bad path to file") //registry is not implemented
		}
	}

	tr.mu.Lock()
	tr.quizzes = quizzes
	tr.mu.Unlock()

	go func() {
		<-ctx.Done()
		tr.cleanup()
	}()

	return nil
}

func (tr *TestRunner) cleanup() {
	time.Sleep(6 * time.Second) //imitating work(actually in-flight Gets saver)

	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.quizzes = nil
}

func (tr *TestRunner) Get(uuid uuid.UUID) (*Quiz, bool) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	if tr.quizzes == nil {
		return nil, false
	}

	q, ok := tr.quizzes[uuid]
	return q, ok
}
