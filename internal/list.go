package internal

import (
	"unsafe"
	_ "unsafe"
)

//var LinkListNodePool = sync.Pool{New: func() interface{} { return new(LinkListNode) }}

type LinkListNode struct {
	next  *LinkListNode
	prev  *LinkListNode
	list  *LinkList
	Value unsafe.Pointer
}

////go:linkname NewLinkListNode github.com/jingyanbin/core/basal.NewLinkListNode
//func NewLinkListNode() *LinkListNode {
//	return LinkListNodePool.Get().(*LinkListNode)
//}

//func (m *LinkListNode) Free() {
//	m.prev = nil
//	m.next = nil
//	m.Value = nil
//	m.list = nil
//	LinkListNodePool.Put(m)
//}

func (m *LinkListNode) Next() *LinkListNode {
	if p := m.next; m.list != nil && p != &m.list.root {
		return p
	}
	return nil
}

func (m *LinkListNode) Prev() *LinkListNode {
	if p := m.prev; m.list != nil && p != &m.list.root {
		return p
	}
	return nil
}

type LinkList struct {
	root LinkListNode
	len  int
}

func (m *LinkList) Init() *LinkList {
	m.root.next = &m.root
	m.root.prev = &m.root
	m.len = 0
	return m
}

//go:linkname NewLinkList github.com/jingyanbin/core/basal.NewLinkList
func NewLinkList() *LinkList { return new(LinkList).Init() }

func (m *LinkList) Len() int { return m.len }

func (m *LinkList) Front() *LinkListNode {
	if m.len == 0 {
		return nil
	}
	return m.root.next
}

func (m *LinkList) Back() *LinkListNode {
	if m.len == 0 {
		return nil
	}
	return m.root.prev
}

func (m *LinkList) lazyInit() {
	if m.root.next == nil {
		m.Init()
	}
}

func (m *LinkList) insert(node, at *LinkListNode) *LinkListNode {
	node.prev = at
	node.next = at.next
	node.prev.next = node
	node.next.prev = node
	node.list = m
	m.len++
	return node
}

func (m *LinkList) insertValue(v unsafe.Pointer, at *LinkListNode) *LinkListNode {
	node := &LinkListNode{Value: v}
	return m.insert(node, at)
}

func (m *LinkList) remove(node *LinkListNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
	node.next = nil
	node.prev = nil
	node.list = nil
	m.len--
}

func (m *LinkList) move(e, at *LinkListNode) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}

func (m *LinkList) Remove(node *LinkListNode) bool {
	if node.list == m {
		m.remove(node)
		return true
	}
	return false
}

func (m *LinkList) PushFront(v unsafe.Pointer) *LinkListNode {
	m.lazyInit()
	return m.insertValue(v, &m.root)
}

func (m *LinkList) PushBack(v unsafe.Pointer) *LinkListNode {
	m.lazyInit()
	return m.insertValue(v, m.root.prev)
}

func (m *LinkList) InsertBefore(v unsafe.Pointer, mark *LinkListNode) *LinkListNode {
	if mark.list != m {
		return nil
	}
	return m.insertValue(v, mark.prev)
}

func (m *LinkList) InsertAfter(v unsafe.Pointer, mark *LinkListNode) *LinkListNode {
	if mark.list != m {
		return nil
	}
	return m.insertValue(v, mark)
}

func (m *LinkList) MoveToFront(node *LinkListNode) {
	if node.list != m || m.root.next == node {
		return
	}
	m.move(node, &m.root)
}

func (m *LinkList) MoveToBack(node *LinkListNode) {
	if node.list != m || m.root.prev == node {
		return
	}
	m.move(node, m.root.prev)
}

func (m *LinkList) MoveBefore(node, mark *LinkListNode) {
	if node.list != m || node == mark || mark.list != m {
		return
	}
	m.move(node, mark.prev)
}

func (m *LinkList) MoveAfter(node, mark *LinkListNode) {
	if node.list != m || node == mark || mark.list != m {
		return
	}
	m.move(node, mark)
}

func (m *LinkList) PushBackList(other *LinkList) {
	m.lazyInit()
	for i, node := other.Len(), other.Front(); i > 0; i, node = i-1, node.Next() {
		m.insertValue(node.Value, m.root.prev)
	}
}

func (m *LinkList) PushFrontList(other *LinkList) {
	m.lazyInit()
	for i, node := other.Len(), other.Back(); i > 0; i, node = i-1, node.Prev() {
		m.insertValue(node.Value, &m.root)
	}
}
