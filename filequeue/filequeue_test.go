package filequeue

import (
	"testing"
)

func BenchmarkFileQueue(b *testing.B) {
	opt := Option{
		MsgFileMaxByte:   MBToByteCount(1),
		PushChanSize:     1000,
		DeletePoppedFile: true,
	}
	q, err := NewFileQueue(opt, nil)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push([]byte("1"))
	}
}

func BenchmarkFileQueuePop(b *testing.B) {
	opt := Option{
		MsgFileMaxByte:   MBToByteCount(1),
		PushChanSize:     1000,
		DeletePoppedFile: true,
	}
	q, err := NewFileQueue(opt, nil)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		q.Push([]byte("1"))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}
