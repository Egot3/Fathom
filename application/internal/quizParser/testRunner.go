package quizparser

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

type TestRunner struct {
	mu         sync.RWMutex
	quizzes    map[uuid.UUID]*Quiz
	generation uint64
	cancel     context.CancelFunc
}

func (tr *TestRunner) Start(ctx context.Context, quizPathsWithUUIDs map[uuid.UUID]string) error {
	quizzes := make(map[uuid.UUID]*Quiz, len(quizPathsWithUUIDs))
	for quizUUID, path := range quizPathsWithUUIDs {
		if !filepath.IsLocal(path) {
			return fmt.Errorf("unsupported path scheme %q: only local paths are currently supported", path) //registry is not implemented
		}
		quiz, err := ParseQuizByPath(path) //no reason to hold the lock when I/O and not writing to tr
		if err != nil {
			return fmt.Errorf("parsing quiz at %q: %w", path, err)
		}
		quizzes[quizUUID] = quiz

	}

	ctx, cancel := context.WithCancel(ctx)

	tr.mu.Lock()
	if tr.cancel != nil {
		tr.cancel()
	}
	tr.cancel = cancel
	tr.quizzes = quizzes
	tr.generation++
	gen := tr.generation
	tr.mu.Unlock()

	go func() {
		<-ctx.Done()
		tr.cleanup(gen)
	}()

	return nil
}

func (tr *TestRunner) cleanup(gen uint64) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	if tr.generation == gen {
		tr.quizzes = nil
	}
}

func (tr *TestRunner) Get(uuid uuid.UUID) (*Quiz, bool) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	q, ok := tr.quizzes[uuid]
	return q, ok
}

func (tr *TestRunner) Stop() {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	if tr.cancel != nil {
		tr.cancel()
		tr.cancel = nil
	}
}
