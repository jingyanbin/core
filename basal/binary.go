package basal

import "encoding/binary"

const zero = byte('0')
const one = byte('1')

//var BytesArrNot16BitMultiple = NewError("[]byte len not is 16bit multiple")
//var BytesArrNot32BitMultiple = NewError("[]byte len not is 32bit multiple")
//var BytesArrNot64BitMultiple = NewError("[]byte len not is 64bit multiple")

var BigEndian = bigEndian{}

type bigEndian struct{}

// byte转二进制byte
func (m *bigEndian) byteToBinary(b byte, buf []byte) {
	if (b>>7)&1 == 1 {
		buf[0] = one
	} else {
		buf[0] = zero
	}
	if (b>>6)&1 == 1 {
		buf[1] = one
	} else {
		buf[1] = zero
	}
	if (b>>5)&1 == 1 {
		buf[2] = one
	} else {
		buf[2] = zero
	}
	if (b>>4)&1 == 1 {
		buf[3] = one
	} else {
		buf[3] = zero
	}
	if (b>>3)&1 == 1 {
		buf[4] = one
	} else {
		buf[4] = zero
	}
	if (b>>2)&1 == 1 {
		buf[5] = one
	} else {
		buf[5] = zero
	}
	if (b>>1)&1 == 1 {
		buf[6] = one
	} else {
		buf[6] = zero
	}
	if (b>>0)&1 == 1 {
		buf[7] = one
	} else {
		buf[7] = zero
	}
}

func (m *bigEndian) BytesToBinary(bs ...byte) []byte {
	buf := make([]byte, 8*len(bs))
	for i, b := range bs {
		m.byteToBinary(b, buf[i*8:])
	}
	return buf
}

func (m *bigEndian) Uint8ToBinary(n ...uint8) []byte {
	return m.BytesToBinary(n...)
}

func (m *bigEndian) Uint16ToBinary(n ...uint16) []byte {
	data := make([]byte, 2*len(n))
	for i, v := range n {
		binary.BigEndian.PutUint16(data[i*2:], v)
	}
	return m.BytesToBinary(data...)
}

func (m *bigEndian) Uint32ToBinary(n ...uint32) []byte {
	data := make([]byte, 4)
	for i, v := range n {
		binary.BigEndian.PutUint32(data[i*4:], v)
	}
	return m.BytesToBinary(data...)
}

func (m *bigEndian) Uint64ToBinary(n ...uint64) []byte {
	data := make([]byte, 8)
	for i, v := range n {
		binary.BigEndian.PutUint64(data[i*8:], v)
	}
	return m.BytesToBinary(data...)
}

func (m *bigEndian) BytesToBinaryString(bs ...byte) string {
	return string(m.BytesToBinary(bs...))
}

func (m *bigEndian) Uint8ToBinaryString(n ...uint8) string {
	return string(m.Uint8ToBinary(n...))
}

func (m *bigEndian) Uint16ToBinaryString(n ...uint16) string {
	return string(m.Uint16ToBinary(n...))
}

func (m *bigEndian) Uint32ToBinaryString(n ...uint32) string {
	return string(m.Uint32ToBinary(n...))
}

func (m *bigEndian) Uint64ToBinaryString(n ...uint64) string {
	return string(m.Uint64ToBinary(n...))
}

//func (*bigEndian) BytesToUint16Arr(bs []byte) ([]uint16, error) {
//	length := len(bs)
//	if length&1 != 0 {
//		return nil, BytesArrNot16BitMultiple
//	}
//	num := length / 2
//	buf := make([]uint16, num)
//	for i := 0; i < length/2; i++ {
//		buf[i] = binary.BigEndian.Uint16(bs[i*2:])
//	}
//	return buf, nil
//}
//
//func (*bigEndian) BytesToUint32Arr(bs []byte) ([]uint32, error) {
//	length := len(bs)
//	if length&3 != 0 {
//		return nil, BytesArrNot32BitMultiple
//	}
//	num := length / 4
//	buf := make([]uint32, num)
//	for i := 0; i < length/4; i++ {
//		buf[i] = binary.BigEndian.Uint32(bs[i*4:])
//	}
//	return buf, nil
//}
//
//func (*bigEndian) BytesToUint64Arr(bs []byte) ([]uint64, error) {
//	length := len(bs)
//	if length&7 != 0 {
//		return nil, BytesArrNot64BitMultiple
//	}
//	num := length / 8
//	buf := make([]uint64, num)
//	for i := 0; i < length/8; i++ {
//		buf[i] = binary.BigEndian.Uint64(bs[i*8:])
//	}
//	return buf, nil
//}
//
//func (*bigEndian) BytesToHex(bs []byte) string {
//	return hex.EncodeToString(bs)
//}
//
//func (*bigEndian) Uint16ArrToBytes(n ...uint16) []byte {
//	buf := make([]byte, 2*len(n))
//	for i, v := range n {
//		binary.BigEndian.PutUint16(buf[2*i:], v)
//	}
//	return buf
//}
//
//func (*bigEndian) Uint32ArrToBytes(n ...uint32) []byte {
//	buf := make([]byte, 4*len(n))
//	for i, v := range n {
//		binary.BigEndian.PutUint32(buf[4*i:], v)
//	}
//	return buf
//}
//
//func (*bigEndian) Uint64ArrToBytes(n ...uint64) []byte {
//	buf := make([]byte, 8*len(n))
//	for i, v := range n {
//		binary.BigEndian.PutUint64(buf[8*i:], v)
//	}
//	return buf
//}

//func (*bigEndian) HexToBytes(s string) ([]byte, error) {
//	return hex.DecodeString(s)
//}
//
//func (*bigEndian) HexToUint16(s string) (uint16, error) {
//	data, err := hex.DecodeString(s)
//	if err != nil {
//		return 0, err
//	}
//	return binary.BigEndian.Uint16(data), nil
//}

//func (*bigEndian) HexToUint32(s string) (uint32, error) {
//	data, err := hex.DecodeString(s)
//	if err != nil {
//		return 0, err
//	}
//	return binary.BigEndian.Uint32(data), nil
//}
//
//func (*bigEndian) HexToUint64(s string) (uint64, error) {
//	data, err := hex.DecodeString(s)
//	if err != nil {
//		return 0, err
//	}
//	return binary.BigEndian.Uint64(data), nil
//}
//
//func (m *bigEndian) HexToUint16Arr(s string) ([]uint16, error) {
//	data, err := hex.DecodeString(s)
//	if err != nil {
//		return nil, err
//	}
//	return m.BytesToUint16Arr(data)
//}
//
//func (m *bigEndian) HexToUint32Arr(s string) ([]uint32, error) {
//	data, err := hex.DecodeString(s)
//	if err != nil {
//		return nil, err
//	}
//	return m.BytesToUint32Arr(data)
//}
//
//func (m *bigEndian) HexToUint64Arr(s string) ([]uint64, error) {
//	data, err := hex.DecodeString(s)
//	if err != nil {
//		return nil, err
//	}
//	return m.BytesToUint64Arr(data)
//}

type littleEndian struct{}

var LittleEndian = littleEndian{}

// byte转二进制byte
func (*littleEndian) byteToBinary(b byte, buf []byte) {
	if (b>>7)&1 == 1 {
		buf[7] = one
	} else {
		buf[7] = zero
	}
	if (b>>6)&1 == 1 {
		buf[6] = one
	} else {
		buf[6] = zero
	}
	if (b>>5)&1 == 1 {
		buf[5] = one
	} else {
		buf[5] = zero
	}
	if (b>>4)&1 == 1 {
		buf[4] = one
	} else {
		buf[4] = zero
	}
	if (b>>3)&1 == 1 {
		buf[3] = one
	} else {
		buf[3] = zero
	}
	if (b>>2)&1 == 1 {
		buf[2] = one
	} else {
		buf[2] = zero
	}
	if (b>>1)&1 == 1 {
		buf[1] = one
	} else {
		buf[1] = zero
	}
	if (b>>0)&1 == 1 {
		buf[0] = one
	} else {
		buf[0] = zero
	}
}

func (m *littleEndian) BytesToBinary(bs ...byte) []byte {
	buf := make([]byte, 8*len(bs))
	for i, b := range bs {
		m.byteToBinary(b, buf[i*8:])
	}
	return buf
}

func (m *littleEndian) Uint8ToBinary(n ...uint8) []byte {
	return m.BytesToBinary(n...)
}

func (m *littleEndian) Uint16ToBinary(n ...uint16) []byte {
	data := make([]byte, 2*len(n))
	for i, v := range n {
		binary.BigEndian.PutUint16(data[i*2:], v)
	}
	return m.BytesToBinary(data...)
}

func (m *littleEndian) Uint32ToBinary(n ...uint32) []byte {
	data := make([]byte, 4)
	for i, v := range n {
		binary.BigEndian.PutUint32(data[i*4:], v)
	}
	return m.BytesToBinary(data...)
}

func (m *littleEndian) Uint64ToBinary(n ...uint64) []byte {
	data := make([]byte, 8)
	for i, v := range n {
		binary.BigEndian.PutUint64(data[i*8:], v)
	}
	return m.BytesToBinary(data...)
}

func (m *littleEndian) BytesToBinaryString(bs ...byte) string {
	return string(m.BytesToBinary(bs...))
}

func (m *littleEndian) Uint8ToBinaryString(n ...uint8) string {
	return string(m.Uint8ToBinary(n...))
}

func (m *littleEndian) Uint16ToBinaryString(n ...uint16) string {
	return string(m.Uint16ToBinary(n...))
}

func (m *littleEndian) Uint32ToBinaryString(n ...uint32) string {
	return string(m.Uint32ToBinary(n...))
}

func (m *littleEndian) Uint64ToBinaryString(n ...uint64) string {
	return string(m.Uint64ToBinary(n...))
}

//func (*littleEndian) BytesToUint16Arr(bs []byte) ([]uint16, error) {
//	length := len(bs)
//	if length&1 != 0 {
//		return nil, BytesArrNot16BitMultiple
//	}
//	num := length / 2
//	buf := make([]uint16, num)
//	for i := 0; i < length/2; i++ {
//		buf[i] = binary.LittleEndian.Uint16(bs[i*2:])
//	}
//	return buf, nil
//}
//
//func (*littleEndian) BytesToUint32Arr(bs []byte) ([]uint32, error) {
//	length := len(bs)
//	if length&3 != 0 {
//		return nil, BytesArrNot32BitMultiple
//	}
//	num := length / 4
//	buf := make([]uint32, num)
//	for i := 0; i < length/4; i++ {
//		buf[i] = binary.LittleEndian.Uint32(bs[i*4:])
//	}
//	return buf, nil
//}
//
//func (*littleEndian) BytesToUint64Arr(bs []byte) ([]uint64, error) {
//	length := len(bs)
//	if length&7 != 0 {
//		return nil, BytesArrNot64BitMultiple
//	}
//	num := length / 8
//	buf := make([]uint64, num)
//	for i := 0; i < length/8; i++ {
//		buf[i] = binary.LittleEndian.Uint64(bs[i*8:])
//	}
//	return buf, nil
//}
//
//func (*littleEndian) BytesToHex(bs []byte) string {
//	return hex.EncodeToString(bs)
//}
//
//func (*littleEndian) Uint16ToBytes(n uint16) []byte {
//	buf := make([]byte, 2)
//	binary.LittleEndian.PutUint16(buf, n)
//	return buf
//}
//
//func (*littleEndian) Uint32ToBytes(n uint32) []byte {
//	buf := make([]byte, 4)
//	binary.LittleEndian.PutUint32(buf, n)
//	return buf
//}
//
//func (*littleEndian) Uint64ToBytes(n uint64) []byte {
//	buf := make([]byte, 8)
//	binary.LittleEndian.PutUint64(buf, n)
//	return buf
//}
//
//func (*littleEndian) Uint16ArrToBytes(n ...uint16) []byte {
//	buf := make([]byte, 2*len(n))
//	for i, v := range n {
//		binary.LittleEndian.PutUint16(buf[2*i:], v)
//	}
//	return buf
//}
//
//func (*littleEndian) Uint32ArrToBytes(n ...uint32) []byte {
//	buf := make([]byte, 4*len(n))
//	for i, v := range n {
//		binary.LittleEndian.PutUint32(buf[4*i:], v)
//	}
//	return buf
//}
//
//func (*littleEndian) Uint64ArrToBytes(n ...uint64) []byte {
//	buf := make([]byte, 8*len(n))
//	for i, v := range n {
//		binary.LittleEndian.PutUint64(buf[8*i:], v)
//	}
//	return buf
//}
//
//func (*littleEndian) HexToBytes(s string) ([]byte, error) {
//	return hex.DecodeString(s)
//}
//
//func (*littleEndian) HexToUint16(s string) (uint16, error) {
//	data, err := hex.DecodeString(s)
//	if err != nil {
//		return 0, err
//	}
//	return binary.LittleEndian.Uint16(data), nil
//}
//
//func (*littleEndian) HexToUint32(s string) (uint32, error) {
//	data, err := hex.DecodeString(s)
//	if err != nil {
//		return 0, err
//	}
//	return binary.LittleEndian.Uint32(data), nil
//}
//
//func (*littleEndian) HexToUint64(s string) (uint64, error) {
//	data, err := hex.DecodeString(s)
//	if err != nil {
//		return 0, err
//	}
//	return binary.LittleEndian.Uint64(data), nil
//}
//
//func (m *littleEndian) HexToUint16Arr(s string) ([]uint16, error) {
//	data, err := hex.DecodeString(s)
//	if err != nil {
//		return nil, err
//	}
//	return m.BytesToUint16Arr(data)
//}
//
//func (m *littleEndian) HexToUint32Arr(s string) ([]uint32, error) {
//	data, err := hex.DecodeString(s)
//	if err != nil {
//		return nil, err
//	}
//	return m.BytesToUint32Arr(data)
//}
//
//func (m *littleEndian) HexToUint64Arr(s string) ([]uint64, error) {
//	data, err := hex.DecodeString(s)
//	if err != nil {
//		return nil, err
//	}
//	return m.BytesToUint64Arr(data)
//}

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
