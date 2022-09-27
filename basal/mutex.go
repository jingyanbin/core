package basal

import (
	internal "github.com/jingyanbin/core/internal"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const checkInterval = time.Second * 10 //检查间隔
const lockedTimeout = 30               //超时

var lockedStackPool = sync.Pool{New: func() interface{} { return new(lockedStack) }}

func getLockedStack() *lockedStack {
	return lockedStackPool.Get().(*lockedStack)
}

type lockedStack struct {
	n     int64
	unix  int64
	stack string
}

func (m *lockedStack) free() {
	lockedStackPool.Put(m)
}

type lockedStacks struct {
	stacks map[int]*lockedStack
	mu     sync.RWMutex
}

func newLockedStacks() *lockedStacks {
	return &lockedStacks{stacks: map[int]*lockedStack{}}
}

func (m *lockedStacks) Get(ptr int) (stack *lockedStack) {
	m.mu.RLock()
	stack = m.stacks[ptr]
	m.mu.RUnlock()
	return
}

func (m *lockedStacks) Set(ptr int, stack *lockedStack) {
	m.mu.Lock()
	m.stacks[ptr] = stack
	m.mu.Unlock()
}

func (m *lockedStacks) Del(ptr int) {
	m.mu.Lock()
	delete(m.stacks, ptr)
	m.mu.Unlock()
}

func (m *lockedStacks) Sub(ptr int) {
	m.mu.Lock()
	stack := m.stacks[ptr]
	if stack != nil {
		v := atomic.AddInt64(&stack.n, -1)
		if v < 1 {
			delete(m.stacks, ptr)
		}
		delete(m.stacks, ptr)
	}
	m.mu.Unlock()
}

func (m *lockedStacks) Range(f func(stack *lockedStack) bool) {
	m.mu.RLock()
	for _, v := range m.stacks {
		if !f(v) {
			break
		}
	}
	m.mu.RUnlock()
}

type lockedMutexManager struct {
	wLockerStacks *lockedStacks
	rLockerStacks *lockedStacks
	mu            sync.RWMutex
}

func (m *lockedMutexManager) wLock(ptr int) {
	if ptr == 0 {
		log.ErrorF("lockedMutexManager wLock error: %v", ptr)
		return
	}
	name, file, line := internal.CallerInFunc(3)
	stack := getLockedStack()
	stack.n = 1
	stack.unix = internal.Unix()
	stack.stack = Sprintf("%s(%s:%d)", name, file, line)
	m.wLockerStacks.Set(ptr, stack)
}

func (m *lockedMutexManager) rLock(ptr int) {
	if ptr == 0 {
		log.ErrorF("lockedMutexManager rLock error: %v", ptr)
		return
	}
	stack := m.rLockerStacks.Get(ptr)
	if stack == nil {
		stack = getLockedStack()
		name, file, line := internal.CallerInFunc(3)
		stack.n = 1
		stack.unix = internal.Unix()
		stack.stack = Sprintf("%s(%s:%d)", name, file, line)
		m.rLockerStacks.Set(ptr, stack)
	} else {
		atomic.AddInt64(&stack.n, 1)
	}
}

func (m *lockedMutexManager) wUnlock(ptr int) {
	if ptr == 0 {
		log.ErrorF("lockedMutexManager wUnlock error: %v", ptr)
		return
	}
	m.wLockerStacks.Del(ptr)
}

func (m *lockedMutexManager) rUnlock(ptr int) {
	if ptr == 0 {
		log.ErrorF("lockedMutexManager rUnlock error: %v", ptr)
		return
	}
	m.rLockerStacks.Sub(ptr)
}

func (m *lockedMutexManager) check() {
	now := internal.Unix()
	m.wLockerStacks.Range(func(stack *lockedStack) bool {
		cha := now - stack.unix
		if cha > lockedTimeout {
			log.ErrorF("please check for deadlock, 长时间未释放写锁(重复加锁,递归加锁,未释放锁等): %s, %v sec, %d", stack.stack, cha, stack.n)
		}
		return true
	})
	m.rLockerStacks.Range(func(stack *lockedStack) bool {
		cha := now - stack.unix
		if cha > lockedTimeout {
			log.ErrorF("please check for deadlock, 长时间未释放读锁(重复加锁,递归加锁,未释放锁等): %s, %v sec, %d", stack.stack, cha, stack.n)
		}
		return true
	})
}

func (m *lockedMutexManager) run() {
	for {
		time.Sleep(checkInterval)
		m.check()
	}
}

func newLockedMutexManager() *lockedMutexManager {
	lockerMgr := &lockedMutexManager{}
	lockerMgr.wLockerStacks = newLockedStacks()
	lockerMgr.rLockerStacks = newLockedStacks()
	go lockerMgr.run()
	return lockerMgr
}

var lockedMus = newLockedMutexManager()

type Mutex struct {
	mu sync.Mutex
}

func (m *Mutex) Lock() {
	m.mu.Lock()
	lockedMus.wLock(*(*int)(unsafe.Pointer(&m)))
}

func (m *Mutex) Unlock() {
	defer m.mu.Unlock()
	lockedMus.wUnlock(*(*int)(unsafe.Pointer(&m)))
}

type rlocker RWMutex

func (r *rlocker) Lock()   { (*RWMutex)(r).RLock() }
func (r *rlocker) Unlock() { (*RWMutex)(r).RUnlock() }

type RWMutex struct {
	rw sync.RWMutex
}

func (m *RWMutex) Lock() {
	m.rw.Lock()
	lockedMus.wLock(*(*int)(unsafe.Pointer(&m)))
}

func (m *RWMutex) Unlock() {
	defer m.rw.Unlock()
	lockedMus.wUnlock(*(*int)(unsafe.Pointer(&m)))
}

func (m *RWMutex) RLock() {
	m.rw.RLock()
	lockedMus.rLock(*(*int)(unsafe.Pointer(&m)))
}

func (m *RWMutex) RUnlock() {
	defer m.rw.RUnlock()
	lockedMus.rUnlock(*(*int)(unsafe.Pointer(&m)))
}

func (m *RWMutex) RLocker() sync.Locker {
	return (*rlocker)(m)
}
