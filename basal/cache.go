package basal

import (
	"github.com/jingyanbin/core/internal"
	"sync"
	"unsafe"
	_ "unsafe"
)

type KVPair[K comparable, V any] struct {
	Key   K
	Value V
}

type LRUCache[K comparable, V any] struct {
	mutex sync.Mutex
	size  int
	list  *internal.LinkList
	cache map[K]*internal.LinkListNode
}

func NewLRUCache[K comparable, V any](size int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		size:  size,
		list:  internal.NewLinkList(),
		cache: make(map[K]*internal.LinkListNode),
	}
}

func (m *LRUCache[K, V]) Len() int {
	return m.list.Len()
}

func (m *LRUCache[K, V]) Size() int {
	return m.size
}

func (m *LRUCache[K, V]) SetSize(size int) {
	m.size = size
}

func (m *LRUCache[K, V]) Get(key K) (v V, ok bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if node, ok2 := m.cache[key]; ok2 {
		m.list.MoveToFront(node)
		kv := (*KVPair[K, V])(node.Value)
		return kv.Value, true
	}
	return v, false
}

func (m *LRUCache[K, V]) Set(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if node, ok := m.cache[key]; ok {
		m.list.MoveToFront(node)
		data := (*KVPair[K, V])(node.Value)
		data.Value = value
		return
	}
	//kv := NewKVPair()
	//kv.Key = unsafe.Pointer(&key)
	//kv.Value = unsafe.Pointer(&value)
	kv := &KVPair[K, V]{Key: key, Value: value}

	nodeNew := m.list.PushFront(unsafe.Pointer(kv))
	m.cache[key] = nodeNew
	if m.list.Len() > m.size {
		back := m.list.Back()
		if back != nil {
			kvBack := (*KVPair[K, V])(back.Value)
			delete(m.cache, kvBack.Key)
			if m.list.Remove(back) {
				//back.Free()
			}
		}
	}
}

func (m *LRUCache[K, V]) Remove(key K) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if node, ok := m.cache[key]; ok {
		delete(m.cache, key)
		if m.list.Remove(node) {
			//node.Free()
		}
	}
}
