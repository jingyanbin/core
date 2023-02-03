package basal

import (
	"github.com/jingyanbin/core/deepcopy"
)

// 有序集合
type SortedSet struct {
	*SortedList
}

// my与b的差集
func (my *SortedSet) Difference(b *SortedSet) *SortedSet {
	c := deepcopy.Copy(my)
	for _, value := range b.buf {
		c.RemoveByKey(my.getKey(value))
	}
	return c
}

// 交集
func (my *SortedSet) Intersection(b *SortedSet) *SortedSet {
	c := NewSortedSet(my.reverse, my.getScore, my.getKey)
	for _, value := range my.buf {
		_, found := b.SearchByKey(my.getKey(value))
		if found {
			c.Add(value)
		}
	}
	return c
}

// 并集
func (my *SortedSet) Union(b *SortedSet) *SortedSet {
	c := deepcopy.Copy(my)
	for _, value := range b.buf {
		c.Add(value)
	}
	return c
}

func (my *SortedSet) Add(v interface{}) bool {
	my.RemoveByKey(my.getKey(v))
	return my.SortedList.Add(v)
}

func NewSortedSet(reverse bool, getScore func(v interface{}) int64, getKey func(v interface{}) int64) *SortedSet {
	return &SortedSet{SortedList: NewSortedList(true, reverse, getScore, getKey)}
}

func NewSortedSetInt(reverse bool) *SortedSet {
	getScore := func(v interface{}) int64 {
		x, err := ToInt64(v)
		if err != nil {
			panic(NewError("NewSortedSetInt min not is int: %v", Type(v)))
		}
		return x
	}
	set := NewSortedSet(
		reverse,
		getScore, getScore)
	return set
}
