package testrunner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"github.com/egot3/fathom/internal/hashutils"
	"github.com/egot3/fathom/internal/quiz"
	quizparser "github.com/egot3/fathom/internal/quizParser"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/samber/lo"
	"github.com/zeebo/xxh3"
)

var (
	ErrQuizNotCached   = errors.New("quiz not cached in runner")
	ErrRunnerInactive  = errors.New("runner not started")
	ErrRunnerPaused    = errors.New("runner is paused")
	ErrRunnerNotPaused = errors.New("runner wasn't paused")
	ErrRunnerExpired   = errors.New("runner expired")
	ErrBadQuizzes      = errors.New("quizUUIDs and quizPathes lens differ")
)

type NotCachedError struct {
	Count int
}

func (e *NotCachedError) Error() string {
	return fmt.Sprintf("%d quizzes not cached in runner", e.Count)
}

func (e *NotCachedError) Is(target error) bool {
	return target == ErrQuizNotCached
}

type concreteTestRunner struct {
	mu         sync.RWMutex
	quizzes    []quiz.Quiz
	generation uint64
	cancel     context.CancelFunc
	timer      *time.Timer
	isPaused   bool
	deadline   time.Time
	pausedAt   time.Time
	TestUUID   uuid.UUID
	GroupUUIDs uuid.UUIDs
	checksum   uint64
}

func NewTestRunner(i do.Injector) (TestRunner, error) {
	return &concreteTestRunner{
		generation: 0,
	}, nil
}

func (tr *concreteTestRunner) CurrentTestUUID() uuid.UUID {
	return tr.TestUUID
}

func (tr *concreteTestRunner) Start(ctx context.Context, duration time.Duration, quizPaths []string, quizUUIDs, groupUUIDs uuid.UUIDs, testUUID uuid.UUID) error {
	if len(quizPaths) != len(quizUUIDs) {
		return ErrBadQuizzes
	}

	quizzes := make([]quiz.Quiz, len(quizPaths))
	for i, path := range quizPaths {
		if !filepath.IsAbs(path) {
			return fmt.Errorf("unsupported path scheme %q: only local paths are currently supported", path) //registry is not implemented
		}
		buf, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		quiz, err := quizparser.ParseQuizByBytes(buf) //no reason to hold the lock when I/O and not writing to tr
		if err != nil {
			return fmt.Errorf("parsing quiz at %q: %w", path, err)
		}
		quiz.UUID = quizUUIDs[i]
		quiz.Checksum = xxh3.Hash(buf)
		quizzes[i] = *quiz
	}

	cctx, cancel := context.WithCancel(ctx)

	tr.mu.Lock()
	if tr.cancel != nil {
		tr.cancel()
	}
	if tr.timer != nil {
		tr.timer.Stop()
	}

	tr.TestUUID = testUUID
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

func (tr *concreteTestRunner) cleanup(gen uint64) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	if tr.generation == gen {
		tr.quizzes = nil
		if tr.cancel != nil {
			tr.cancel()
			tr.cancel = nil
		}
		tr.TestUUID = uuid.Nil
	}
}

func (tr *concreteTestRunner) Get(quizUUID uuid.UUID) (*quiz.Quiz, error) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	if tr.cancel == nil {
		return nil, ErrRunnerInactive
	}

	q, ok := lo.Find(tr.quizzes, func(quiz quiz.Quiz) bool {
		return quiz.UUID == quizUUID
	})
	if !ok {
		return nil, ErrQuizNotCached
	}
	return &q, nil
}

func (tr *concreteTestRunner) Stop() {
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
func (tr *concreteTestRunner) UpsertQuiz(quizPaths []string, quizUUIDs uuid.UUIDs) error {
	if len(quizPaths) != len(quizUUIDs) {
		return ErrBadQuizzes
	}

	tr.mu.RLock()
	if tr.cancel == nil { // reading a bunch of files might fry the potato
		tr.mu.RUnlock()
		return ErrRunnerInactive
	}
	tr.mu.RUnlock()

	quizzes := make([]quiz.Quiz, len(quizPaths))
	for i, path := range quizPaths {
		if !filepath.IsAbs(path) {
			return fmt.Errorf("unsupported path scheme %q: only abs paths are currently supported", path) //registry is not implemented
		}
		buf, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		quiz, err := quizparser.ParseQuizByBytes(buf) //no reason to hold the lock when I/O and not writing to tr
		if err != nil {
			return fmt.Errorf("parsing quiz at %q: %w", path, err)
		}
		quiz.UUID = quizUUIDs[i]
		quizzes[i] = *quiz

	}

	tr.mu.Lock()
	defer tr.mu.Unlock()
	if tr.cancel == nil { //TOCTOU
		return ErrRunnerInactive
	}
	tr.quizzes = append(tr.quizzes, quizzes...)

	ultimateChecksum := make([]uint64, len(tr.quizzes))
	lo.ForEach(tr.quizzes, func(quiz quiz.Quiz, _ int) {
		ultimateChecksum = append(ultimateChecksum, quiz.Checksum)
	})

	tr.checksum = hashutils.HashHashes(ultimateChecksum)

	return nil
}

// remember: parioal deletion, return 200 207 OR 404
func (tr *concreteTestRunner) RemoveQuiz(uuids uuid.UUIDs) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	if tr.cancel == nil {
		return ErrRunnerInactive
	}

	oldL := len(tr.quizzes)
	tr.quizzes = lo.Filter(tr.quizzes, func(quiz quiz.Quiz, _ int) bool {
		return !slices.Contains(uuids, quiz.UUID)
	})
	nf := len(uuids) - (oldL - len(tr.quizzes))

	if nf != 0 {
		return &NotCachedError{Count: nf}
	}

	return nil
}

func (tr *concreteTestRunner) ExtendTime(duration time.Duration) error {
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
func (tr *concreteTestRunner) Pause() error {
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

func (tr *concreteTestRunner) Resume() error {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	if tr.cancel == nil {
		return ErrRunnerInactive
	}
	if !tr.isPaused {
		return ErrRunnerNotPaused
	}

	sincePause := time.Since(tr.pausedAt)
	tr.deadline = tr.deadline.Add(sincePause)

	if sincePause <= 0 {
		tr.cancel()
		tr.cancel = nil
		tr.quizzes = nil
		tr.isPaused = false
		return ErrRunnerInactive
	}

	if time.Until(tr.deadline) <= 0 {
		tr.cancel()
		tr.cancel = nil
		tr.quizzes = nil
		tr.isPaused = false
		return ErrRunnerExpired
	}

	tr.timer.Reset(sincePause)

	/* log.Printf("paused for: %v, now %v, new deadline: %v",
	sincePause.String(),
	time.Now().Format(time.TimeOnly),
	tr.deadline.Format(time.TimeOnly)) */

	tr.isPaused = false
	return nil
}

func (tr *concreteTestRunner) Deadline() (*time.Time, error) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	if tr.cancel == nil {
		return nil, ErrRunnerInactive
	}

	return &tr.deadline, nil
}

func (tr *concreteTestRunner) AllowedGroupUUIDs() uuid.UUIDs {
	return tr.GroupUUIDs
}

func (tr *concreteTestRunner) GetAll() uuid.UUIDs {
	tr.mu.RLock()
	quizzesLocal := tr.quizzes
	tr.mu.RUnlock()

	quizzes := lo.Map(quizzesLocal, func(quiz quiz.Quiz, _ int) uuid.UUID {
		return quiz.UUID
	})

	return quizzes
}

func (tr *concreteTestRunner) Checksum() uint64 {
	return tr.checksum
}
