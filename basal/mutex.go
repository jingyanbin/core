package basal

import (
	internal "github.com/jingyanbin/core/internal"
	"runtime"
	"sync"
	"time"
)

const lockTimeout = 30000 //超时

type Mutex struct {
	mu         sync.Mutex
	lockedTime int64
}

func (m *Mutex) Lock() {
	start := internal.UnixMs()
	var now int64
	var waitMs int64
	for !m.mu.TryLock() {
		now = internal.UnixMs()
		waitMs = now - start
		if waitMs < 33 {
			runtime.Gosched()
		} else {
			if m.lockedTime > 0 {
				if cha := now - m.lockedTime; cha > lockTimeout {
					name, file, line := internal.CallerInFunc(2)
					stack := Sprintf("%s(%s:%d)", name, file, line)
					log.ErrorF("please check for deadlock Mutex Lock, 长时间未释放写锁(重复加锁,递归加锁,未释放锁等): %s, %ds", stack, cha/1000)
					time.Sleep(time.Second)
					continue
				}
			}
			time.Sleep(time.Millisecond)
		}
	}
	m.lockedTime = internal.UnixMs()
}

func (m *Mutex) Unlock() {
	m.lockedTime = 0
	m.mu.Unlock()
}

func (m *Mutex) Exec(f func()) {
	m.Lock()
	defer m.Unlock()
	f()
}

type rLocker RWMutex

func (r *rLocker) Lock()   { (*RWMutex)(r).RLock() }
func (r *rLocker) Unlock() { (*RWMutex)(r).RUnlock() }

type RWMutex struct {
	rw         sync.RWMutex
	lockedTime int64
}

func (m *RWMutex) Lock() {
	start := internal.UnixMs()
	var now int64
	var waitMs int64
	for !m.rw.TryLock() {
		now = internal.UnixMs()
		waitMs = now - start
		if waitMs < 33 {
			runtime.Gosched()
		} else {
			if m.lockedTime > 0 {
				if cha := now - m.lockedTime; cha > lockTimeout {
					name, file, line := internal.CallerInFunc(2)
					stack := Sprintf("%s(%s:%d)", name, file, line)
					log.ErrorF("please check for deadlock RWMutex Lock, 长时间未释放写锁(重复加锁,递归加锁,未释放锁等): %s, %ds", stack, cha/1000)
					time.Sleep(time.Second)
					continue
				}
			}
			time.Sleep(time.Millisecond)
		}
	}
	m.lockedTime = internal.UnixMs()
}

func (m *RWMutex) Unlock() {
	m.lockedTime = 0
	m.rw.Unlock()
}

func (m *RWMutex) RLock() {
	start := internal.UnixMs()
	var now int64
	var waitMs int64
	for !m.rw.TryRLock() {
		now = internal.UnixMs()
		waitMs = now - start
		if waitMs < 33 {
			runtime.Gosched()
		} else {
			if m.lockedTime > 0 {
				if cha := now - m.lockedTime; cha > lockTimeout {
					name, file, line := internal.CallerInFunc(2)
					stack := Sprintf("%s(%s:%d)", name, file, line)
					log.ErrorF("please check for deadlock RWMutex RLock, 长时间未释放写锁(重复加锁,递归加锁,未释放锁等): %s, %ds", stack, cha/1000)
					time.Sleep(time.Second)
					continue
				}
			}
			time.Sleep(time.Millisecond)
		}
	}
	m.lockedTime = internal.UnixMs()
}

func (m *RWMutex) RUnlock() {
	m.lockedTime = 0
	m.rw.RUnlock()
}

func (m *RWMutex) RLocker() sync.Locker {
	return (*rLocker)(m)
}

func (m *RWMutex) Exec(f func()) {
	m.Lock()
	defer m.Unlock()
	f()
}

func (m *RWMutex) RExec(f func()) {
	m.RLock()
	defer m.RUnlock()
	f()
}
