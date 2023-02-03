package basal

import (
	"github.com/jingyanbin/core/log"
	"strings"
	"sync"
	"time"
	"unsafe"
	_ "unsafe"
)

type LRUCaches[K comparable, V any] struct {
	shard        int
	slot         func(K) int
	caches       []*LRUCache[K, V]
	printInfoSec int32
}

func (m *LRUCaches[K, V]) run() {
	ticker := time.NewTicker(time.Second)
	infoArr := make([]string, m.shard)
	var last int64
	for {
		select {
		case now := <-ticker.C:
			if int32(now.Unix()-last) >= m.printInfoSec {
				last = now.Unix()
				for i, v := range m.caches {
					infoArr[i] = Sprintf("%d/%d", v.Len(), v.list.Len())
				}
				log.Info("LRUCaches: %s", strings.Join(infoArr, ","))
			}
		}
	}
}

func (m *LRUCaches[K, V]) SetSize(size int) {
	for _, v := range m.caches {
		v.SetSize(size)
	}
}

func (m *LRUCaches[K, V]) Get(key K, new func(K) V) (v V, ok bool) {
	return m.caches[m.slot(key)].Get(key, new)
}

// NewLRUCaches[K comparable, V any]
//
//	@Description:
//	@param size  每个节点存储数据大小
//	@param shard 分片大小
//	@param printInfoSec 打印信息间隔秒数
//	@param slot 计算槽位的方法,返回槽位的index
//	@return *LRUCaches[K,V] 返回缓存
func NewLRUCaches[K comparable, V any](size int, shard int, printInfoSec int32, slot func(K) int) *LRUCaches[K, V] {
	caches := &LRUCaches[K, V]{
		shard:        shard,
		slot:         slot,
		caches:       make([]*LRUCache[K, V], shard),
		printInfoSec: printInfoSec,
	}
	for i := 0; i < shard; i++ {
		caches.caches[i] = NewLRUCache[K, V](size)
	}
	if printInfoSec > 0 {
		go caches.run()
	}
	return caches
}

type KVPair[K comparable, V any] struct {
	Key   K
	Value V
}

type LRUCache[K comparable, V any] struct {
	mu          sync.RWMutex
	size        int
	list        *LinkList
	moveFrontCh chan *LinkListNode
	pushFrontCh chan *LinkListNode
	removeCh    chan *LinkListNode
	cache       map[K]*LinkListNode
}

func NewLRUCache[K comparable, V any](size int) *LRUCache[K, V] {
	cache := &LRUCache[K, V]{
		size:        size,
		list:        NewLinkList(),
		cache:       make(map[K]*LinkListNode, size),
		moveFrontCh: make(chan *LinkListNode, size),
		pushFrontCh: make(chan *LinkListNode, size),
		removeCh:    make(chan *LinkListNode, size),
	}
	go cache.run()
	return cache
}

func (m *LRUCache[K, V]) Len() int {
	return len(m.cache)
}

func (m *LRUCache[K, V]) Size() int {
	return m.size
}

func (m *LRUCache[K, V]) SetSize(size int) {
	m.size = size
}

func (m *LRUCache[K, V]) Get(key K, new func(K) V) (v V, ok bool) {
	m.mu.RLock()
	node, ok1 := m.cache[key]
	m.mu.RUnlock()
	if ok1 {
		kv := (*KVPair[K, V])(node.Value)
		m.moveFrontCh <- node
		return kv.Value, true
	}
	if new == nil {
		return v, false
	}
	m.mu.Lock()
	node, ok1 = m.cache[key]
	if ok1 {
		m.mu.Unlock()
		kv := (*KVPair[K, V])(node.Value)
		m.moveFrontCh <- node
		return kv.Value, true
	}

	v = new(key)
	kv := &KVPair[K, V]{Key: key, Value: v}
	node = &LinkListNode{Value: unsafe.Pointer(kv)}
	m.cache[key] = node
	m.mu.Unlock()
	m.pushFrontCh <- node

	return v, true

}

func (m *LRUCache[K, V]) removeExpired() {
	if dLen := m.list.Len() - m.size; dLen > 0 {
		//st := time.Now()
		m.mu.Lock()
		for i := 0; i < dLen; i++ {
			back := m.list.Back()
			if back == nil {
				break
			}
			kvBack := (*KVPair[K, V])(back.Value)
			delete(m.cache, kvBack.Key)
			m.list.Remove(back)
		}
		m.mu.Unlock()
		//et := time.Now()
		//log.Info("2==================%v, %vms, %v, %v", dLen, et.Sub(st).Milliseconds(), et.Sub(st).Microseconds(), et.Sub(st).Nanoseconds())
	}
}

func (m *LRUCache[K, V]) run() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case node := <-m.moveFrontCh:
			m.list.MoveToFront(node)
		case node := <-m.pushFrontCh:
			m.list.PushNodeFront(node)
		case node := <-m.removeCh:
			m.list.Remove(node)
		case <-ticker.C:
			m.removeExpired()
		}
	}
}

func (m *LRUCache[K, V]) Remove(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if node, ok := m.cache[key]; ok {
		delete(m.cache, key)
		m.removeCh <- node
	}
}
