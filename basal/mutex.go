package basal

import (
	internal "github.com/jingyanbin/core/internal"
	"sync"
	"time"
	"unsafe"
)

type lockedStack struct {
	unix  int64
	stack string
}

type lockedMutexManager struct {
	ptrLockers map[int]lockedStack
	mu         sync.RWMutex
}

func (m *lockedMutexManager) add(ptr int) {
	if ptr == 0 {
		log.ErrorF("lockedMutexManager add error: %v", ptr)
		return
	}
	name, file, line := internal.CallerInFunc(3)
	m.mu.Lock()
	m.ptrLockers[ptr] = lockedStack{unix: internal.Unix(), stack: Sprintf("%s(%s:%d)", name, file, line)}
	m.mu.Unlock()
}

func (m *lockedMutexManager) del(ptr int) {
	if ptr == 0 {
		log.ErrorF("lockedMutexManager del error: %v", ptr)
		return
	}
	m.mu.Lock()
	delete(m.ptrLockers, ptr)
	m.mu.Unlock()
}

func (m *lockedMutexManager) check() {
	now := internal.Unix()
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, v := range m.ptrLockers {
		cha := now - v.unix
		if cha > 30 {
			log.ErrorF("please check for deadlock,重复加锁,递归加锁,未释放锁等: %s, %v sec", v.stack, cha)
		}
	}
}

func (m *lockedMutexManager) run() {
	for {
		time.Sleep(time.Second * 10)
		m.check()
	}
}

func newLockedMutexManager() *lockedMutexManager {
	lockerMgr := &lockedMutexManager{ptrLockers: map[int]lockedStack{}}
	go lockerMgr.run()
	return lockerMgr
}

var lockedMus = newLockedMutexManager()

type Mutex struct {
	sync.Mutex
}

func (m *Mutex) Lock() {
	m.Mutex.Lock()
	lockedMus.add(*(*int)(unsafe.Pointer(&m)))
}

func (m *Mutex) Unlock() {
	defer m.Mutex.Unlock()
	lockedMus.del(*(*int)(unsafe.Pointer(&m)))
}

type RWMutex struct {
	sync.RWMutex
}

func (m *RWMutex) Lock() {
	m.RWMutex.Lock()
	lockedMus.add(*(*int)(unsafe.Pointer(&m)))
}

func (m *RWMutex) Unlock() {
	defer m.RWMutex.Unlock()
	lockedMus.del(*(*int)(unsafe.Pointer(&m)))
}
