package basal

import (
	"container/list"
	"sync"
)

type pair struct {
	k interface{}
	v interface{}
}

type LRUCache struct {
	mutex sync.Mutex
	size  int
	list  *list.List
	cache map[interface{}]*list.Element
}

func NewLRUCache(size int) *LRUCache {
	return &LRUCache{
		size:  size,
		list:  list.New(),
		cache: make(map[interface{}]*list.Element),
	}
}

func (my *LRUCache) Get(key interface{}) (interface{}, bool) {
	my.mutex.Lock()
	defer my.mutex.Unlock()
	if elem, ok := my.cache[key]; ok {
		my.list.MoveToFront(elem)
		return elem.Value.(*pair).v, true
	}
	return nil, false
}

func (my *LRUCache) Set(key interface{}, value interface{}) {
	my.mutex.Lock()
	defer my.mutex.Unlock()
	if elem, ok := my.cache[key]; ok {
		my.list.MoveToFront(elem)
		elem.Value = &pair{k: key, v: value}
	} else {
		elemNew := my.list.PushFront(&pair{k: key, v: value})
		my.cache[key] = elemNew
		if my.list.Len() >= my.size {
			back := my.list.Back()
			if back != nil {
				delete(my.cache, back.Value.(*pair).k)
				my.list.Remove(back)
			}
		}
	}
}

func (my *LRUCache) Remove(key interface{}) {
	my.mutex.Lock()
	defer my.mutex.Unlock()
	if elem, ok := my.cache[key]; ok {
		delete(my.cache, key)
		my.list.Remove(elem)
	}
}

func (my *LRUCache) Len() int {
	my.mutex.Lock()
	defer my.mutex.Unlock()
	return my.list.Len()
}

func (my *LRUCache) Size() int {
	my.mutex.Lock()
	defer my.mutex.Unlock()
	return my.size
}

func (my *LRUCache) SetSize(size int) {
	my.mutex.Lock()
	defer my.mutex.Unlock()
	my.size = size
}
