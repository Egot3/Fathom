package quizparser

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"sync"
	"time"
)

var (
	ErrQuizNotCached   = errors.New("quiz not cached in runner")
	ErrRunnerInactive  = errors.New("runner not started")
	ErrRunnerPaused    = errors.New("runner is paused")
	ErrRunnerNotPaused = errors.New("runner wasn't paused")
	ErrRunnerExpired   = errors.New("runner expired")
)

type NotCachedError struct {
	Ids []int
}

func (e *NotCachedError) Error() string {
	return fmt.Sprintf("%d quizzes not cached in runner", len(e.Ids))
}

func (e *NotCachedError) Is(target error) bool {
	return target == ErrQuizNotCached
}

type TestRunner struct {
	mu         sync.RWMutex
	quizzes    []*Quiz
	generation uint64
	cancel     context.CancelFunc
	timer      *time.Timer
	isPaused   bool
	deadline   time.Time
	pausedAt   time.Time
}

// Ctx must explicitly hold the lifetime of test
func (tr *TestRunner) Start(ctx context.Context, duration time.Duration, quizPaths []string) error {
	quizzes := make([]*Quiz, len(quizPaths))
	for i, path := range quizPaths {
		if !filepath.IsLocal(path) {
			return fmt.Errorf("unsupported path scheme %q: only local paths are currently supported", path) //registry is not implemented
		}
		quiz, err := ParseQuizByPath(path) //no reason to hold the lock when I/O and not writing to tr
		if err != nil {
			return fmt.Errorf("parsing quiz at %q: %w", path, err)
		}
		quizzes[i] = quiz

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

func (tr *TestRunner) Get(id int) (*Quiz, error) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	if tr.cancel == nil {
		return nil, ErrRunnerInactive
	}

	if ok := len(tr.quizzes) > id; !ok {
		return nil, ErrQuizNotCached
	}
	q := tr.quizzes[id]

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

// acquire path quiz pairs from db via uuid
// called upsert as maps.Copy overwrites dest on collision
func (tr *TestRunner) UpsertQuiz(quizPaths []string) error {
	tr.mu.RLock()
	if tr.cancel == nil { // reading a bunch of files might fry the potato
		tr.mu.RUnlock()
		return ErrRunnerInactive
	}
	tr.mu.RUnlock()

	quizzes := make([]*Quiz, len(quizPaths))
	for i, path := range quizPaths {
		if !filepath.IsLocal(path) {
			return fmt.Errorf("unsupported path scheme %q: only local paths are currently supported", path) //registry is not implemented
		}
		quiz, err := ParseQuizByPath(path) //no reason to hold the lock when I/O and not writing to tr
		if err != nil {
			return fmt.Errorf("parsing quiz at %q: %w", path, err)
		}
		quizzes[i] = quiz

	}

	tr.mu.Lock()
	defer tr.mu.Unlock()
	if tr.cancel == nil { //TOCTOU
		return ErrRunnerInactive
	}
	tr.quizzes = slices.Clone(quizzes)

	return nil
}

// remember: parioal deletion, return 200 207 OR 404
func (tr *TestRunner) RemoveQuiz(ids []int) error {
	tr.mu.RLock()
	if tr.cancel == nil {
		tr.mu.RUnlock()
		return ErrRunnerInactive
	}
	tr.mu.RUnlock() //just to keep up with upsert

	notFound := make([]int, 0)
	tr.mu.Lock()
	if tr.cancel == nil {
		return ErrRunnerInactive
	}
	for _, id := range ids {
		if ok := len(tr.quizzes) > id; !ok {
			notFound = append(notFound, id)
			continue
		}

		tr.quizzes = append(tr.quizzes[:id], tr.quizzes[id+1:]...)
	}
	tr.mu.Unlock()

	if len(notFound) != 0 {
		return &NotCachedError{Ids: notFound}
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
