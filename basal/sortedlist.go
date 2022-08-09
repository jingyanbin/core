package basal

//有序列表
type SortedList struct {
	buf     []interface{}
	Cmp     func(min, max interface{}) int //返回值 min < max return 1, min == max return 0, min > max return -1
	Reverse bool                           //反序
}

func (my *SortedList) Len() int {
	return len(my.buf)
}

func (my *SortedList) Slice() []interface{} {
	return my.buf
}

func (my *SortedList) binarySearch(value interface{}) (index int, found bool) {
	length := len(my.buf)
	if length == 0 {
		return 0, false
	}
	start := 0
	end := length - 1
	var cmp int
	for {
		index = start + (end-start)/2
		cmp = my.Cmp(value, my.buf[index])
		if my.Reverse {
			if cmp == 1 {
				start = index + 1
			} else if cmp == -1 {
				end = index - 1
			} else if cmp == 0 {
				return index, true
			} else {
				panic(NewError("func cmp rturn value not in (-1, 0, 1)"))
			}
		} else { //正序
			if cmp == 1 {
				end = index - 1
			} else if cmp == -1 {
				start = index + 1
			} else if cmp == 0 {
				return index, true
			} else {
				panic(NewError("func cmp rturn value not in (-1, 0, 1)"))
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

func (my *SortedList) Del(index int) bool {
	if index < 0 || index >= len(my.buf) {
		return false
	}
	my.buf = append(my.buf[:index], my.buf[index+1:]...)
	return true
}

func (my *SortedList) Find(v interface{}) (index int, found bool) {
	return my.binarySearch(v)
}

func (my *SortedList) add(v interface{}, repeat bool) bool {
	index, found := my.binarySearch(v)
	if found && repeat == false {
		return false
	}
	my.buf = append(my.buf, v)
	copy(my.buf[index+1:], my.buf[index:])
	my.buf[index] = v
	return true
}

func (my *SortedList) Add(v interface{}) bool {
	return my.add(v, true)
}

func (my *SortedList) Remove(v interface{}) bool {
	index, found := my.binarySearch(v)
	if found {
		my.buf = append(my.buf[:index], my.buf[index+1:]...)
		return true
	} else {
		return false
	}
}

func (my *SortedList) Clear() {
	my.buf = my.buf[:0]
}
