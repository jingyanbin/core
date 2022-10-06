package internal

import (
	"sync"
	"sync/atomic"
)

type OnceSuccess struct {
	m    sync.Mutex
	done uint32
}

func (m *OnceSuccess) Success() bool {
	return atomic.LoadUint32(&m.done) == 1
}

func (m *OnceSuccess) Do(f func() bool) bool {
	if atomic.LoadUint32(&m.done) == 1 {
		return true
	}
	m.m.Lock()
	defer m.m.Unlock()
	if m.done == 0 {
		if !f() {
			return false
		}
		atomic.StoreUint32(&m.done, 1)
	}
	return true
}

func (m *OnceSuccess) DoError(f func() error) error {
	if atomic.LoadUint32(&m.done) == 1 {
		return nil
	}
	m.m.Lock()
	defer m.m.Unlock()
	if m.done == 0 {
		if err := f(); err != nil {
			return err
		}
		atomic.StoreUint32(&m.done, 1)
	}
	return nil
}
