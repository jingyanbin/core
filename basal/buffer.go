package basal

import (
	"github.com/jingyanbin/core/internal"
)

var BufferPool = internal.BufferPool

type Buffer = internal.Buffer

func NewBuffer(size int) *Buffer {
	return internal.NewBuffer(size)
}
