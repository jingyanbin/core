package basal

import "github.com/jingyanbin/core/deepcopy"

//有序集合
type SortedSet struct {
	SortedList
}

func (my *SortedSet) Add(v interface{}) bool {
	return my.add(v, false)
}

//并集
func (my *SortedSet) Union(b *SortedSet) *SortedSet {
	c := deepcopy.Copy(my).(*SortedSet)
	for _, value := range b.buf {
		c.Add(value)
	}
	return c
}

//my与b的差集
func (my *SortedSet) Difference(b *SortedList) *SortedList {
	c := deepcopy.Copy(my).(*SortedList)
	for _, value := range b.buf {
		c.Remove(value)
	}
	return c
}

//交集
func (my *SortedSet) Intersection(b *SortedSet) *SortedSet {
	c := NewSortedSet(my.Cmp, my.Reverse)
	for _, value := range my.buf {
		_, found := b.binarySearch(value)
		if found {
			c.Add(value)
		}
	}
	return c
}

func NewSortedSet(cmp func(min, max interface{}) int, reverse bool) *SortedSet {
	return &SortedSet{SortedList{Cmp: cmp, Reverse: reverse}}
}

func NewSortedSetInt(reverse bool) *SortedSet {
	set := &SortedSet{}
	set.Reverse = reverse
	set.Cmp = func(min, max interface{}) int {
		x, err := ToInt64(min)
		if err != nil {
			panic(NewError("SortedSet min not is int: %v", min))
		}
		y, err := ToInt64(max)
		if err != nil {
			panic(NewError("SortedSet max not is int: %v", min))
		}
		if x < y {
			return 1
		} else if x > y {
			return -1
		} else {
			return 0
		}
	}
	return set
}
