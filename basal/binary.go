package basal

import (
	"encoding/binary"
	"encoding/hex"
)

const zero = byte('0')
const one = byte('1')

var BytesArrNot16BitMultiple = NewError("[]byte len not is 16bit multiple")
var BytesArrNot32BitMultiple = NewError("[]byte len not is 32bit multiple")
var BytesArrNot64BitMultiple = NewError("[]byte len not is 64bit multiple")

//byte转二进制byte
func ByteToBinary(b byte) []byte {
	buf := make([]byte, 8)
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
	return buf
}

//byte转二进制字符串
func ByteToBinaryString(b byte) string {
	return string(ByteToBinary(b))
}

//bytes转二进制bytes
func BytesToBinary(bs []byte) []byte {
	res := make([]byte, 0, len(bs)*8)
	for _, b := range bs {
		res = append(res, ByteToBinary(b)...)
	}
	return res
}

//bytes转二进制字符串
func BytesToBinaryString(bs []byte) string {
	return string(BytesToBinary(bs))
}

//按位设值
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

//按位设值
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

//按位获取值
func GetBitUint64(v uint64, offset int) bool {
	if offset > 63 || offset < 0 {
		return false
	}
	return ((v >> offset) & 1) == 1
}

//按位获取值
func GetBitUint32(v uint32, offset int) bool {
	if offset > 31 || offset < 0 {
		return false
	}
	return ((v >> offset) & 1) == 1
}

type bigEndian struct{}

func (*bigEndian) Uint16ToBinaryString(v uint16) string {
	bs := BigEndian.Uint16ToBytes(v)
	buf := make([]byte, 0, 16)
	buf = append(buf, ByteToBinary(bs[0])...)
	buf = append(buf, ByteToBinary(bs[1])...)
	return string(buf)
}

func (*bigEndian) Uint32ToBinaryString(v uint32) string {
	bs := BigEndian.Uint32ToBytes(v)
	buf := make([]byte, 0, 32)
	buf = append(buf, ByteToBinary(bs[0])...)
	buf = append(buf, ByteToBinary(bs[1])...)
	buf = append(buf, ByteToBinary(bs[2])...)
	buf = append(buf, ByteToBinary(bs[3])...)
	return string(buf)
}

func (*bigEndian) Uint64ToBinaryString(v uint64) string {
	bs := BigEndian.Uint64ToBytes(v)
	buf := make([]byte, 0, 64)
	buf = append(buf, ByteToBinary(bs[0])...)
	buf = append(buf, ByteToBinary(bs[1])...)
	buf = append(buf, ByteToBinary(bs[2])...)
	buf = append(buf, ByteToBinary(bs[3])...)
	buf = append(buf, ByteToBinary(bs[4])...)
	buf = append(buf, ByteToBinary(bs[5])...)
	buf = append(buf, ByteToBinary(bs[6])...)
	buf = append(buf, ByteToBinary(bs[7])...)
	return string(buf)
}

func (*bigEndian) BytesToUint16Arr(bs []byte) ([]uint16, error) {
	length := len(bs)
	if length&1 != 0 {
		return nil, BytesArrNot16BitMultiple
	}
	num := length / 2
	buf := make([]uint16, num)
	for i := 0; i < length/2; i++ {
		buf[i] = binary.BigEndian.Uint16(bs[i*2:])
	}
	return buf, nil
}

func (*bigEndian) BytesToUint32Arr(bs []byte) ([]uint32, error) {
	length := len(bs)
	if length&3 != 0 {
		return nil, BytesArrNot32BitMultiple
	}
	num := length / 4
	buf := make([]uint32, num)
	for i := 0; i < length/4; i++ {
		buf[i] = binary.BigEndian.Uint32(bs[i*4:])
	}
	return buf, nil
}

func (*bigEndian) BytesToUint64Arr(bs []byte) ([]uint64, error) {
	length := len(bs)
	if length&7 != 0 {
		return nil, BytesArrNot64BitMultiple
	}
	num := length / 8
	buf := make([]uint64, num)
	for i := 0; i < length/8; i++ {
		buf[i] = binary.BigEndian.Uint64(bs[i*8:])
	}
	return buf, nil
}

func (*bigEndian) BytesToHex(bs []byte) string {
	return hex.EncodeToString(bs)
}

func (*bigEndian) Uint16ToBytes(n uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, n)
	return buf
}

func (*bigEndian) Uint32ToBytes(n uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, n)
	return buf
}

func (*bigEndian) Uint64ToBytes(n uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, n)
	return buf
}

func (*bigEndian) Uint16ArrToBytes(n ...uint16) []byte {
	buf := make([]byte, 2*len(n))
	for i, v := range n {
		binary.BigEndian.PutUint16(buf[2*i:], v)
	}
	return buf
}

func (*bigEndian) Uint32ArrToBytes(n ...uint32) []byte {
	buf := make([]byte, 4*len(n))
	for i, v := range n {
		binary.BigEndian.PutUint32(buf[4*i:], v)
	}
	return buf
}

func (*bigEndian) Uint64ArrToBytes(n ...uint64) []byte {
	buf := make([]byte, 8*len(n))
	for i, v := range n {
		binary.BigEndian.PutUint64(buf[8*i:], v)
	}
	return buf
}

func (*bigEndian) HexToBytes(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func (*bigEndian) HexToUint16(s string) (uint16, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(data), nil
}

func (*bigEndian) HexToUint32(s string) (uint32, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(data), nil
}

func (*bigEndian) HexToUint64(s string) (uint64, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(data), nil
}

func (m *bigEndian) HexToUint16Arr(s string) ([]uint16, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return m.BytesToUint16Arr(data)
}

func (m *bigEndian) HexToUint32Arr(s string) ([]uint32, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return m.BytesToUint32Arr(data)
}

func (m *bigEndian) HexToUint64Arr(s string) ([]uint64, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return m.BytesToUint64Arr(data)
}

var BigEndian = &bigEndian{}

//littleEndian
type littleEndian struct{}

func (*littleEndian) Uint16ToBinaryString(v uint16) string {
	bs := BigEndian.Uint16ToBytes(v)
	buf := make([]byte, 0, 16)
	buf = append(buf, ByteToBinary(bs[0])...)
	buf = append(buf, ByteToBinary(bs[1])...)
	return string(buf)
}

func (*littleEndian) Uint32ToBinaryString(v uint32) string {
	bs := BigEndian.Uint32ToBytes(v)
	buf := make([]byte, 0, 32)
	buf = append(buf, ByteToBinary(bs[0])...)
	buf = append(buf, ByteToBinary(bs[1])...)
	buf = append(buf, ByteToBinary(bs[2])...)
	buf = append(buf, ByteToBinary(bs[3])...)
	return string(buf)
}

func (*littleEndian) Uint64ToBinaryString(v uint64) string {
	bs := BigEndian.Uint64ToBytes(v)
	buf := make([]byte, 0, 64)
	buf = append(buf, ByteToBinary(bs[0])...)
	buf = append(buf, ByteToBinary(bs[1])...)
	buf = append(buf, ByteToBinary(bs[2])...)
	buf = append(buf, ByteToBinary(bs[3])...)
	buf = append(buf, ByteToBinary(bs[4])...)
	buf = append(buf, ByteToBinary(bs[5])...)
	buf = append(buf, ByteToBinary(bs[6])...)
	buf = append(buf, ByteToBinary(bs[7])...)
	return string(buf)
}

func (*littleEndian) BytesToUint16Arr(bs []byte) ([]uint16, error) {
	length := len(bs)
	if length&1 != 0 {
		return nil, BytesArrNot16BitMultiple
	}
	num := length / 2
	buf := make([]uint16, num)
	for i := 0; i < length/2; i++ {
		buf[i] = binary.LittleEndian.Uint16(bs[i*2:])
	}
	return buf, nil
}

func (*littleEndian) BytesToUint32Arr(bs []byte) ([]uint32, error) {
	length := len(bs)
	if length&3 != 0 {
		return nil, BytesArrNot32BitMultiple
	}
	num := length / 4
	buf := make([]uint32, num)
	for i := 0; i < length/4; i++ {
		buf[i] = binary.LittleEndian.Uint32(bs[i*4:])
	}
	return buf, nil
}

func (*littleEndian) BytesToUint64Arr(bs []byte) ([]uint64, error) {
	length := len(bs)
	if length&7 != 0 {
		return nil, BytesArrNot64BitMultiple
	}
	num := length / 8
	buf := make([]uint64, num)
	for i := 0; i < length/8; i++ {
		buf[i] = binary.LittleEndian.Uint64(bs[i*8:])
	}
	return buf, nil
}

func (*littleEndian) BytesToHex(bs []byte) string {
	return hex.EncodeToString(bs)
}

func (*littleEndian) Uint16ToBytes(n uint16) []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, n)
	return buf
}

func (*littleEndian) Uint32ToBytes(n uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, n)
	return buf
}

func (*littleEndian) Uint64ToBytes(n uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, n)
	return buf
}

func (*littleEndian) Uint16ArrToBytes(n ...uint16) []byte {
	buf := make([]byte, 2*len(n))
	for i, v := range n {
		binary.LittleEndian.PutUint16(buf[2*i:], v)
	}
	return buf
}

func (*littleEndian) Uint32ArrToBytes(n ...uint32) []byte {
	buf := make([]byte, 4*len(n))
	for i, v := range n {
		binary.LittleEndian.PutUint32(buf[4*i:], v)
	}
	return buf
}

func (*littleEndian) Uint64ArrToBytes(n ...uint64) []byte {
	buf := make([]byte, 8*len(n))
	for i, v := range n {
		binary.LittleEndian.PutUint64(buf[8*i:], v)
	}
	return buf
}

func (*littleEndian) HexToBytes(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func (*littleEndian) HexToUint16(s string) (uint16, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(data), nil
}

func (*littleEndian) HexToUint32(s string) (uint32, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(data), nil
}

func (*littleEndian) HexToUint64(s string) (uint64, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(data), nil
}

func (m *littleEndian) HexToUint16Arr(s string) ([]uint16, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return m.BytesToUint16Arr(data)
}

func (m *littleEndian) HexToUint32Arr(s string) ([]uint32, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return m.BytesToUint32Arr(data)
}

func (m *littleEndian) HexToUint64Arr(s string) ([]uint64, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return m.BytesToUint64Arr(data)
}

var LittleEndian = &littleEndian{}
