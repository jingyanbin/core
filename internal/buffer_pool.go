package internal

import (
	"sync"
)

const bufferPoolNumber = 10
const minBufferSize = 256

func makeBuffer(size int) *Buffer {
	buf := make(Buffer, 0, size)
	return &buf
}

type bufferPool struct {
	sync.Pool
}

func (pool *bufferPool) init(index int) {
	pool.New = func() interface{} {
		if index > 0 {
			return makeBuffer(index * minBufferSize)
		} else {
			return makeBuffer(minBufferSize / 2)
		}
	}
}

type bufferPools [bufferPoolNumber]*bufferPool

func (pools *bufferPools) Init() {
	for i := 0; i < bufferPoolNumber; i++ {
		pools[i] = &bufferPool{}
		pools[i].init(i)
	}
}

func (pools *bufferPools) New(size int) *Buffer {
	index := size / minBufferSize
	if index < 0 || index >= bufferPoolNumber {
		return makeBuffer(size)
	} else {
		return pools[index].Get().(*Buffer)
	}
}

func (pools *bufferPools) Free(buf *Buffer) {
	index := cap(*buf) / minBufferSize
	if index >= 0 && index < bufferPoolNumber {
		*buf = (*buf)[:0]
		pools[index].Put(buf)
	}
}

var buffersMgr bufferPools

func init() {
	buffersMgr.Init()
}
