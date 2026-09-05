package utils

import (
	"sync"
	"sync/atomic"
)

type LockManager struct {
	mu    sync.Mutex
	locks map[string]*lockEntry
}

type lockEntry struct {
	mu   sync.Mutex
	refs int32
}

func NewLockManager() *LockManager {
	return &LockManager{locks: make(map[string]*lockEntry)}
}

func (lm *LockManager) Lock(key string) func() {
	lm.mu.Lock()
	e := lm.locks[key]
	if e == nil {
		e = &lockEntry{}
		lm.locks[key] = e
	}
	atomic.AddInt32(&e.refs, 1)
	lm.mu.Unlock()

	e.mu.Lock()

	var once sync.Once
	return func() {
		once.Do(func() {
			e.mu.Unlock()

			if atomic.AddInt32(&e.refs, -1) == 0 {
				lm.mu.Lock()
				if lm.locks[key] == e && atomic.LoadInt32(&e.refs) == 0 {
					delete(lm.locks, key)
				}
				lm.mu.Unlock()
			}
		})
	}
}
