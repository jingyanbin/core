package basal

import "sync"

import _ "unsafe"

type Map[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

func (m *Map[K, V]) Set(key K, value V) (old V) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[K]V)
	}
	v := m.data[key]
	m.data[key] = value
	m.mu.Unlock()
	return v
}

func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	m.mu.RLock()
	value, ok = m.data[key]
	m.mu.RUnlock()
	return
}

// Delete
//
//	@Description: 函数内不可再调本Map方法
//	@receiver m
//	@param key
//	@param f
//	@return bool
func (m *Map[K, V]) Delete(key K, f func(value V) bool) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		return false
	}
	v, ok := m.data[key]
	if !ok {
		return false
	}
	del := true
	if f != nil {
		del = f(v)
	}
	if del {
		delete(m.data, key)
	}
	return del
}

// Range
//
//	@Description: 函数内不可再调本Map方法
//	@receiver m
//	@param f
//	@return bool
func (m *Map[K, V]) Range(f func(key K, value V) bool) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if !f(k, v) {
			return false
		}
	}
	return true
}

func (m *Map[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}
