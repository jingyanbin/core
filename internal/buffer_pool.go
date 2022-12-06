package internal

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const bufferPoolNumber = 10
const bufferPoolIndexMax = bufferPoolNumber - 1

const minBufferSize = 1 << 10

var bufferLevel = 1
var bufferStatistics bool

type bufferPool struct {
	sync.Pool
	size int64

	freeTotal uint64 //释放总内存
	freeCount uint32 //释放次数
	newCount  uint32 //申请次数
}

func (m *bufferPool) Get() *Buffer {
	if bufferStatistics {
		atomic.AddUint32(&m.newCount, 1)
	}
	return m.Pool.Get().(*Buffer)
}

func (m *bufferPool) Put(buf *Buffer) {
	if bufferStatistics {
		atomic.AddUint64(&m.freeTotal, uint64(buf.Len()))
		atomic.AddUint32(&m.freeCount, 1)
	}
	buf.Clear()
	m.Pool.Put(buf)
}

func (m *bufferPool) init(index int) {
	m.size = int64((index + 1) * minBufferSize * bufferLevel)
	//log.InfoF("================init: index: %v, size: %v", index, m.size)
	m.New = func() interface{} {
		if index > 0 {
			return &Buffer{buf: make([]byte, 0, m.size), index: uint8(index)}
		} else {
			return &Buffer{buf: make([]byte, 0, m.size), index: uint8(index)}
		}
	}
}

type bufferPools [bufferPoolNumber]*bufferPool

func (m *bufferPools) init() {
	for i := 0; i < bufferPoolNumber; i++ {
		m[i] = &bufferPool{}
		m[i].init(i)
	}
}

func (m *bufferPools) SetLevel(level int) {
	bufferLevel = level
}

func (m *bufferPools) SetStatistics(statistics bool) {
	bufferStatistics = statistics
}

func (m *bufferPools) Info() string {
	var s string
	for i, pool := range m {
		var avg uint64
		if pool.freeCount > 0 {
			avg = pool.freeTotal / uint64(pool.freeCount)
		}
		s += fmt.Sprintf("\nid: %d, free total: %d, free/new: %d/%d, avg: %d", i, pool.freeTotal, pool.freeCount, pool.newCount, avg)
	}
	return s
}

func (m *bufferPools) new(size int) *Buffer {
	bufferSize := minBufferSize * bufferLevel
	index := size / bufferSize
	if size%bufferSize == 0 {
		index--
	}
	//log.InfoF("================new: %v", index)
	if index < 0 {
		return m[0].Get()
	} else if index < bufferPoolNumber {
		return m[index].Get()
	} else {
		return &Buffer{buf: make([]byte, 0, size), index: uint8(index)}
	}
}

func (m *bufferPools) free(buf *Buffer) {
	index := buf.index
	if index < bufferPoolNumber {
		m[index].Put(buf)
	} else {
		m[bufferPoolIndexMax].Put(buf)
	}
}

var BufferPool bufferPools

func init() {
	BufferPool.init()
}
