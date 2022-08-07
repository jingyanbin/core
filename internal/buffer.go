package internal

import (
	"strconv"
	_ "unsafe"
)

type Buffer []byte

func (buf *Buffer) AppendByte(b byte) {
	*buf = append(*buf, b)
}

func (buf *Buffer) AppendBytes(bs ...byte) {
	*buf = append(*buf, bs...)
}

func (buf *Buffer) AppendString(s string) {
	*buf = append(*buf, s...)
}

func (buf *Buffer) AppendStrings(ss ...string) {
	for _, s := range ss {
		*buf = append(*buf, s...)
	}
}

func (buf *Buffer) AppendInt(n, w int) {
	ItoA((*[]byte)(buf), n, w)
}

func (buf *Buffer) AppendFloat(f float64) {
	*buf = strconv.AppendFloat(*buf, f, 'f', -1, 64)
}

func (buf *Buffer) Bytes() []byte {
	return *buf
}

func (buf *Buffer) ToString() string {
	return string(*buf)
}

func (buf *Buffer) Clear() {
	*buf = (*buf)[:0]
}

func (buf *Buffer) Cap() int {
	return cap(*buf)
}

func (buf *Buffer) Free() {
	buffersMgr.Free(buf)
}

//go:linkname NewBuffer github.com/jingyanbin/core/basal.NewBuffer
func NewBuffer(size int) *Buffer {
	return buffersMgr.New(size)
}
