package basal

import (
	"fmt"
)

// 有序列表
type SortedList struct {
	buf         []interface{}
	scoreRepeat bool                      //排序值是否可重复
	reverse     bool                      //反序
	getScore    func(v interface{}) int64 //获取分数函数
	getKey      func(v interface{}) int64 //获取key int
}

func (my *SortedList) String() string {
	return fmt.Sprintf("%v", my.buf)
}

func (my *SortedList) reduceSpace() {
	length := len(my.buf)
	max := cap(my.buf)
	if length == 0 {
		my.buf = nil
	} else {
		curMax := max / 2
		if curMax > length {
			buf := make([]interface{}, length, curMax)
			copy(buf, my.buf)
			my.buf = buf
		}
	}
}

func (my *SortedList) Len() int {
	return len(my.buf)
}

func (my *SortedList) Cap() int {
	return cap(my.buf)
}

func (my *SortedList) Slice() []interface{} {
	return my.buf
}

func (my *SortedList) cmp(min, max int64) int {
	if min < max {
		return 1
	} else if min > max {
		return -1
	} else {
		return 0
	}
}

func (my *SortedList) SearchByKey(key int64) (int, bool) {
	for idx, item := range my.buf {
		if my.getKey(item) == key {
			return idx, true
		}
	}
	return 0, false
}

func (my *SortedList) SearchByScore(score int64) (index int, found bool) {
	length := len(my.buf)
	if length == 0 {
		return 0, false
	}
	start := 0
	end := length - 1
	var cmp int
	for {
		index = start + (end-start)/2
		cmp = my.cmp(score, my.getScore(my.buf[index]))
		if my.reverse {
			if cmp == 1 {
				start = index + 1
			} else if cmp == -1 {
				end = index - 1
			} else if cmp == 0 {
				return index, true
			} else {
				panic(NewError("func cmp return value not in (-1, 0, 1)"))
			}
		} else { //正序
			if cmp == 1 {
				end = index - 1
			} else if cmp == -1 {
				start = index + 1
			} else if cmp == 0 {
				return index, true
			} else {
				panic(NewError("func cmp return value not in (-1, 0, 1)"))
			}
		}
		if start > end {
			return start, false
		}
	}
}

func (my *SortedList) Front() (v interface{}, found bool) {
	if len(my.buf) > 0 {
		return my.buf[0], true
	} else {
		return nil, false
	}
}

func (my *SortedList) Back() (v interface{}, found bool) {
	length := len(my.buf)
	if length > 0 {
		return my.buf[length-1], true
	} else {
		return nil, false
	}
}

func (my *SortedList) PopFront() (v interface{}, found bool) {
	if len(my.buf) > 0 {
		v = my.buf[0]
		my.buf = my.buf[1:]
		my.reduceSpace()
		return v, true
	} else {
		return nil, false
	}
}

func (my *SortedList) PopBack() (v interface{}, found bool) {
	length := len(my.buf)
	if length > 0 {
		v = my.buf[length-1]
		my.buf = my.buf[:length-1]
		my.reduceSpace()
		return v, true
	} else {
		return nil, false
	}
}

func (my *SortedList) Get(index int) (v interface{}, found bool) {
	length := len(my.buf)
	if length > 0 {
		if index < 0 || index >= length {
			return nil, false
		}
		return my.buf[index], true
	} else {
		return nil, false
	}
}

func (my *SortedList) Add(v interface{}) bool {
	index, found := my.SearchByScore(my.getScore(v))
	if found && my.scoreRepeat == false {
		return false
	}
	my.buf = append(my.buf, v)
	copy(my.buf[index+1:], my.buf[index:])
	my.buf[index] = v
	return true
}

func (my *SortedList) RemoveByIndex(index int) bool {
	if index < 0 || index >= len(my.buf) {
		return false
	}
	my.buf = append(my.buf[:index], my.buf[index+1:]...)
	my.reduceSpace()
	return true
}

func (my *SortedList) RemoveByScore(score int64) bool {
	index, found := my.SearchByScore(score)
	if found {
		my.buf = append(my.buf[:index], my.buf[index+1:]...)
		my.reduceSpace()
		return true
	} else {
		return false
	}
}

func (my *SortedList) RemoveByKey(key int64) bool {
	index, found := my.SearchByKey(key)
	if found {
		my.buf = append(my.buf[:index], my.buf[index+1:]...)
		my.reduceSpace()
		return true
	} else {
		return false
	}
}

//func (my *SortedList) Reset() {
//	my.buf = my.buf[:0]
//}

func (my *SortedList) Clear() {
	my.buf = nil
}

func NewSortedList(scoreRepeat, reverse bool, getScore func(v interface{}) int64, getKey func(v interface{}) int64) *SortedList {
	return &SortedList{scoreRepeat: scoreRepeat, reverse: reverse, getScore: getScore, getKey: getKey}
}
