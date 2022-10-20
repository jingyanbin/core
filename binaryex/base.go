package binaryex

func SetBitByte(v byte, t bool, offset int) byte {
	if offset > 7 || offset < 0 {
		return v
	}
	if t {
		v = v | (1 << offset)
	} else {
		v = v &^ (1 << offset)
	}
	return v
}

func SetBitUint8(v uint8, t bool, offset int) uint8 {
	if offset > 7 || offset < 0 {
		return v
	}
	if t {
		v = v | (1 << offset)
	} else {
		v = v &^ (1 << offset)
	}
	return v
}

func SetBitUint16(v uint16, t bool, offset int) uint16 {
	if offset > 15 || offset < 0 {
		return v
	}
	if t {
		v = v | (1 << offset)
	} else {
		v = v &^ (1 << offset)
	}
	return v
}
func SetBitUint32(v uint32, t bool, offset int) uint32 {
	if offset > 31 || offset < 0 {
		return v
	}
	if t {
		v = v | (1 << offset)
	} else {
		v = v &^ (1 << offset)
	}
	return v
}

func SetBitUint64(v uint64, t bool, offset int) uint64 {
	if offset > 63 || offset < 0 {
		return v
	}
	if t {
		v = v | (1 << offset)
	} else {
		v = v &^ (1 << offset)
	}
	return v
}

func GetBitByte(v byte, offset int) bool {
	if offset > 7 || offset < 0 {
		return false
	}
	return ((v >> offset) & 1) == 1
}

func GetBitUint8(v uint8, offset int) bool {
	if offset > 7 || offset < 0 {
		return false
	}
	return ((v >> offset) & 1) == 1
}

func GetBitUint16(v uint16, offset int) bool {
	if offset > 15 || offset < 0 {
		return false
	}
	return ((v >> offset) & 1) == 1
}

func GetBitUint32(v uint32, offset int) bool {
	if offset > 31 || offset < 0 {
		return false
	}
	return ((v >> offset) & 1) == 1
}

func GetBitUint64(v uint64, offset int) bool {
	if offset > 63 || offset < 0 {
		return false
	}
	return ((v >> offset) & 1) == 1
}

type BytesBinary []byte

func (m *BytesBinary) expansion(offset int) {
	n := offset/8 + 1
	cha := n - len(*m)
	if cha < 0 {
		return
	}
	size := cha + len(*m)
	buf := make([]byte, size, size)
	copy(buf, *m)
	*m = buf
}

func (m *BytesBinary) ReSize(size int) {
	if size > 0 {
		maxIndex := size - 1
		n := maxIndex/8 + 1
		if n == len(*m) {
			return
		}
		buf := make([]byte, n, n)
		copy(buf, *m)
		*m = buf
	} else {
		*m = nil
	}
}

func (m *BytesBinary) Clear() {
	for i := range *m {
		(*m)[i] = 0
	}
}

func (m *BytesBinary) Set(t bool, offset int) {
	m.expansion(offset)
	index := offset / 8
	bitOffset := offset % 8
	(*m)[index] = SetBitByte((*m)[index], t, bitOffset)
}

func (m *BytesBinary) Get(offset int) bool {
	index := offset / 8
	bitOffset := offset % 8
	return GetBitByte((*m)[index], bitOffset)
}

func (m *BytesBinary) GetBinaryString() string {
	return LittleEndian.BytesToBinaryString(*m...)
}
