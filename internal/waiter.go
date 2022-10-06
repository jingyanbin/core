package internal

import "time"

type Waiter struct {
	ch     chan struct{}
	init   OnceSuccess
	closed OnceSuccess
}

func (m *Waiter) toInit() bool {
	m.ch = make(chan struct{}, 0)
	return true
}

func (m *Waiter) toClose() bool {
	close(m.ch)
	return true
}

func (m *Waiter) Done() {
	m.init.Do(m.toInit)
	m.closed.Do(m.toClose)
}

func (m *Waiter) Wait(timeout time.Duration) {
	m.init.Do(m.toInit)
	if timeout > 0 {
		ticker := time.NewTicker(timeout)
		defer ticker.Stop()
		select {
		case <-m.ch:
		case <-ticker.C:
		}
	} else {
		select {
		case <-m.ch:
		}
	}
}
