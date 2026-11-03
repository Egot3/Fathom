package testrunner

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/egot3/fathom/internal/hashutils"
	"github.com/egot3/fathom/internal/logging"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/zeebo/xxh3"
)

var ErrManagerClosed = errors.New("manager is shutting down")

type Manager struct {
	// noCopy brothers
	mu sync.RWMutex
	wg sync.WaitGroup

	closed    bool
	runners   map[uint64]TestRunner
	totalizer Totalizer
}

func (m *Manager) Start(ctx context.Context, duration time.Duration,
	quizPaths []string, quizUUIDs, groupUUIDs uuid.UUIDs, testUUID uuid.UUID,
) (TestRunner, error) {
	key := deriveKey(groupUUIDs, testUUID)

	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return nil, ErrManagerClosed
	}
	if _, exists := m.runners[key]; exists {
		m.mu.Unlock()
		return nil, ErrAlreadyRunning
	}
	tr := &concreteTestRunner{}
	m.runners[key] = tr
	m.wg.Add(1)
	m.mu.Unlock()

	if err := tr.start(ctx, duration, quizPaths, quizUUIDs, groupUUIDs, testUUID); err != nil {
		m.remove(key, tr)
		m.wg.Done()
		return nil, err
	}

	m.mu.RLock()
	closed := m.closed
	m.mu.RUnlock()

	if closed {
		tr.Stop()
		m.remove(key, tr)
		m.wg.Done()
		return nil, ErrManagerClosed
	}

	go func() {
		defer m.wg.Done()
		m.supervise(context.WithoutCancel(ctx), key, tr)
	}() // whenthe beheaded context

	return tr, nil
}

func (m *Manager) remove(key uint64, tr *concreteTestRunner) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if cur, ok := m.runners[key]; ok && cur == tr {
		delete(m.runners, key)
	}
}

func (m *Manager) Get(key uint64) (TestRunner, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	tr, ok := m.runners[key]
	return tr, ok
}

// deterministic key derivement
func deriveKey(groupUUIDs uuid.UUIDs, testUUID uuid.UUID) uint64 {
	grStrings := groupUUIDs.Strings()
	slices.SortFunc(grStrings, strings.Compare)

	checksums := make([]uint64, len(groupUUIDs)+1)
	for i, str := range grStrings {
		checksums[i] = xxh3.HashString(str)
	}
	checksums[len(checksums)-1] = xxh3.Hash(testUUID[:])
	return hashutils.HashHashes(checksums)
}

func (m *Manager) GetAll() []uint64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return slices.Collect(maps.Keys(m.runners))
}

func (m *Manager) IsQuizRunning(searchedUUID uuid.UUID) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, runner := range m.runners {
		if slices.Contains(runner.Quizzes(), searchedUUID) {
			return true
		}
	}
	return false
}

func (m *Manager) AllRunners() map[uint64]uuid.UUID {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tests := make(map[uint64]uuid.UUID, len(m.runners))
	for key, r := range m.runners {
		tests[key] = r.Test()
	}

	return tests
}

func NewManagerWith(t Totalizer) *Manager {
	return &Manager{runners: map[uint64]TestRunner{}, totalizer: t}
}

func NewManager(i do.Injector) (*Manager, error) {
	logger := do.MustInvoke[*slog.Logger](i)
	t, err := do.InvokeAs[Totalizer](i)
	if err != nil {
		logger.Warn("no Totalizer registered; expired runners will NOT be totalized")
		t = nopTotalizer{}
	}
	return NewManagerWith(t), nil
}

func (m *Manager) supervise(ctx context.Context, key uint64, tr *concreteTestRunner) {
	<-tr.done
	defer m.remove(key, tr)

	if tr.reason == ReasonStopped {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	logger := logging.LoggerFromContext(ctx).With(slog.String("layer", "manager"), slog.Uint64("runner_key", key))
	if err := m.totalizer.TotalizeAll(ctx, tr.Test(), tr.Groups(), tr.StartedAt()); err != nil {
		logger.Error("couldn't totalize expired runner", slog.String("Error", err.Error()))
	}
}

func (m *Manager) ShutdownWithContext(ctx context.Context) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Started shutting down the manager")

	m.mu.Lock()
	m.closed = true
	live := slices.Collect(maps.Values(m.runners))
	m.mu.Unlock()

	for _, tr := range live {
		tr.stopFor(ReasonShutdowned)
	}

	c := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(c)
	}()

	select {
	case <-ctx.Done():
		logger.Error("Couldn't shut down manager", slog.String("Error", ctx.Err().Error()))
		return fmt.Errorf("Couldn't shut down manager: %w", ctx.Err())
	case <-c:
		logger.Info("Shutdown of manager complete")
		return nil
	}
}
