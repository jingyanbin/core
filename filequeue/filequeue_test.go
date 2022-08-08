package filequeue

import (
	"testing"
)

func BenchmarkFileQueue(b *testing.B) {
	opt := &Options{
		ConfDataDir:      "file_queue",
		MsgFileDir:       "testlog",
		FileNamePrefix:   "test",
		MsgFileMaxByte:   MBToByteCount(1),
		PushChanSize:     1000,
		DeletePoppedFile: true,
	}
	q, err := NewFileQueue(opt)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push([]byte("1"))
	}
}

func BenchmarkFileQueuePop(b *testing.B) {
	opt := &Options{
		ConfDataDir:      "file_queue",
		MsgFileDir:       "testlog",
		FileNamePrefix:   "test",
		MsgFileMaxByte:   MBToByteCount(1),
		PushChanSize:     1000,
		DeletePoppedFile: true,
		ReadCount:        false,
	}
	q, err := NewFileQueue(opt)
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
