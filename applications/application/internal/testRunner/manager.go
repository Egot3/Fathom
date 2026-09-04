package testrunner

import (
	"context"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/egot3/fathom/internal/hashutils"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/zeebo/xxh3"
)

type runnerInfo struct {
	TestUUID   uuid.UUID
	GroupUUIDs uuid.UUIDs
}

type Manager struct {
	mu         sync.RWMutex
	runners    map[uint64]TestRunner
	runnerInfo map[uint64]runnerInfo
}

func (m *Manager) Start(ctx context.Context, duration time.Duration,
	quizPaths []string, quizUUIDs, groupUUIDs uuid.UUIDs, testUUID uuid.UUID,
) (TestRunner, error) {
	key := deriveKey(groupUUIDs, testUUID)

	m.mu.Lock()
	if _, exists := m.runners[key]; exists {
		m.mu.Unlock()
		return nil, ErrAlreadyRunning
	}
	tr := &concreteTestRunner{}
	m.runners[key] = tr
	m.runnerInfo[key] = runnerInfo{
		TestUUID: testUUID, GroupUUIDs: groupUUIDs,
	}
	m.mu.Unlock()

	if err := tr.start(ctx, duration, quizPaths, quizUUIDs,
		func() { m.remove(key, tr) }); err != nil {
		m.remove(key, tr)
		return nil, err
	}
	return tr, nil
}

func (m *Manager) remove(key uint64, tr *concreteTestRunner) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if cur, ok := m.runners[key]; ok && cur == tr {
		delete(m.runners, key)
		delete(m.runnerInfo, key)
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

	keys := make([]uint64, len(m.runners))
	i := 0
	for key, _ := range m.runners {
		keys[i] = key
		i++
	}

	return keys
}

func (m *Manager) GetInfo(key uint64) (runnerInfo, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ri, ok := m.runnerInfo[key]
	return ri, ok
}

func NewManager(i do.Injector) (Manager, error) {
	return Manager{
		mu:      sync.RWMutex{},
		runners: map[uint64]TestRunner{},
	}, nil
}
