package basal

import (
	"sync/atomic"
	"time"
)

type Waiter struct {
	ch     chan struct{}
	init   OnceSuccess
	closed OnceSuccess
	n      int32
}

func (m *Waiter) toInit() bool {
	m.ch = make(chan struct{}, 0)
	return true
}

func (m *Waiter) toClose() bool {
	close(m.ch)
	return true
}

func (m *Waiter) toAdd(n int32) int32 {
	return atomic.AddInt32(&m.n, n)
}

func (m *Waiter) Add(n uint32) {
	m.toAdd(int32(n))
}

func (m *Waiter) Done() {
	m.init.Do(m.toInit)
	if m.toAdd(-1) < 1 {
		m.closed.Do(m.toClose)
	}
}

// 返回false表示超时
func (m *Waiter) Wait(timeout time.Duration) bool {
	m.init.Do(m.toInit)
	if timeout > 0 {
		ticker := time.NewTimer(timeout)
		defer ticker.Stop()
		select {
		case <-m.ch:
			return true
		case <-ticker.C:
			return false
		}
	} else {
		select {
		case <-m.ch:
			return true
		}
	}
}
