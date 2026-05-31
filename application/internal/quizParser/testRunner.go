package quizparser

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrQuizNotCached   = errors.New("quiz not cached in runner")
	ErrRunnerInactive  = errors.New("runner not started")
	ErrRunnerPaused    = errors.New("runner is paused")
	ErrRunnerNotPaused = errors.New("runner wasn't paused")
	ErrRunnerExpired   = errors.New("runner expired")
)

type NotCachedError struct {
	Uuids uuid.UUIDs
}

func (e *NotCachedError) Error() string {
	return fmt.Sprintf("%d quizzes not cached in runner", len(e.Uuids))
}

func (e *NotCachedError) Is(target error) bool {
	return target == ErrQuizNotCached
}

type TestRunner struct {
	mu         sync.RWMutex
	quizzes    map[uuid.UUID]*Quiz
	generation uint64
	cancel     context.CancelFunc
	timer      *time.Timer
	isPaused   bool
	deadline   time.Time
	pausedAt   time.Time
}

// Ctx must explicitly hold the lifetime of test
func (tr *TestRunner) Start(ctx context.Context, duration time.Duration, quizPathsWithUUIDs map[uuid.UUID]string) error {
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

	cctx, cancel := context.WithCancel(ctx)

	tr.mu.Lock()
	if tr.cancel != nil {
		tr.cancel()
	}
	if tr.timer != nil {
		tr.timer.Stop()
	}

	tr.cancel = cancel
	tr.quizzes = quizzes
	tr.isPaused = false
	tr.deadline = time.Now().Add(duration)

	tr.generation++
	gen := tr.generation

	tr.timer = time.AfterFunc(time.Until(tr.deadline), func() {
		tr.cleanup(gen)
	})
	localTimer := tr.timer
	tr.mu.Unlock()

	go func() {

		<-cctx.Done()
		if localTimer.Stop() {
			tr.cleanup(gen)
		}

	}()

	return nil
}

func (tr *TestRunner) cleanup(gen uint64) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	if tr.generation == gen {
		tr.quizzes = nil
		if tr.cancel != nil {
			tr.cancel()
			tr.cancel = nil
		}
	}
}

func (tr *TestRunner) Get(uuid uuid.UUID) (*Quiz, error) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	if tr.cancel == nil {
		return nil, ErrRunnerInactive
	}
	q, ok := tr.quizzes[uuid]
	if !ok {
		return nil, ErrQuizNotCached
	}

	return q, nil
}

func (tr *TestRunner) Stop() {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	if tr.cancel != nil {
		tr.cancel()
		tr.cancel = nil
	}
	if tr.timer != nil {
		tr.timer.Stop()
	}
	tr.quizzes = nil
}

// acquire uuid-path quiz pairs from db via uuid
// called upsert as maps.Copy overwrites dest on collision
func (tr *TestRunner) UpsertQuiz(quizPathsWithUUIDs map[uuid.UUID]string) error {
	tr.mu.RLock()
	if tr.cancel == nil { // reading a bunch of files might fry the potato
		tr.mu.RUnlock()
		return ErrRunnerInactive
	}
	tr.mu.RUnlock()

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

	tr.mu.Lock()
	defer tr.mu.Unlock()
	if tr.cancel == nil { //TOCTOU
		return ErrRunnerInactive
	}
	maps.Copy(tr.quizzes, quizzes)

	return nil
}

// remember: parioal deletion, return 200 207 OR 404
func (tr *TestRunner) RemoveQuiz(quizUUIDs uuid.UUIDs) error {
	tr.mu.RLock()
	if tr.cancel == nil {
		tr.mu.RUnlock()
		return ErrRunnerInactive
	}
	tr.mu.RUnlock() //just to keep up with upsert

	notFound := make(uuid.UUIDs, 0)
	tr.mu.Lock()
	if tr.cancel == nil {
		return ErrRunnerInactive
	}
	for _, quizUUID := range quizUUIDs {
		if _, ok := tr.quizzes[quizUUID]; !ok {
			notFound = append(notFound, quizUUID)
			continue
		}
		delete(tr.quizzes, quizUUID)
	}
	tr.mu.Unlock()

	if len(notFound) != 0 {
		return &NotCachedError{Uuids: notFound}
	}

	return nil
}

func (tr *TestRunner) ExtendTime(duration time.Duration) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	if tr.cancel == nil {
		return ErrRunnerInactive
	}

	tr.deadline = tr.deadline.Add(duration)
	if !tr.isPaused {
		if !tr.timer.Stop() {
			tr.generation++
		}
		tr.timer.Reset(time.Until(tr.deadline))
	}

	return nil
}

// doesn't lock the runner for the whole pause letting dialer know what's up
func (tr *TestRunner) Pause() error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	if tr.cancel == nil {
		return ErrRunnerInactive
	}
	if tr.isPaused {
		return ErrRunnerPaused
	}
	if !tr.timer.Stop() {
		return ErrRunnerExpired //expiration
	}

	tr.isPaused = true
	tr.pausedAt = time.Now()
	return nil
}

func (tr *TestRunner) Resume() error {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	if tr.cancel == nil {
		return ErrRunnerInactive
	}
	if !tr.isPaused {
		return ErrRunnerNotPaused
	}

	tr.deadline = tr.deadline.Add(time.Since(tr.pausedAt))
	remaining := time.Until(tr.deadline)
	if remaining <= 0 {
		tr.cancel()
		tr.cancel = nil
		tr.quizzes = nil
		tr.isPaused = false
		return ErrRunnerInactive
	}

	tr.timer.Reset(remaining)
	tr.isPaused = false
	return nil
}
