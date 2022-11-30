package uuid

import (
	"github.com/jingyanbin/core/internal"
	"testing"
)

func TestUUID(t *testing.T) {
	opt := Option{}
	opt.IndexBits = 18
	opt.WorkerIdBits = 10
	opt.TimeBits = 35
	opt.TimeLatest = true
	opt.WorkerId = 1
	opt.Epoch = internal.UnixMs() + (3600 * 24 * 397 * 1000)
	gen := NewGenerator(opt)
	for i := 0; i < 10; i++ {
		uid := gen.UUID()
		ms, index, serverId := gen.DeUUID(uid)
		t.Logf("ms: %v, index: %v, serverId: %v", ms, index, serverId)
	}
}

func BenchmarkUUID(b *testing.B) {
	opt := Option{}
	opt.IndexBits = 18
	opt.WorkerIdBits = 10
	opt.TimeBits = 35
	opt.TimeLatest = false
	opt.WorkerId = 1
	opt.Epoch = internal.UnixMs() + (3600 * 24 * 397 * 1000)
	var gen IGenerator = NewGenerator(opt)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gen.UUID()
	}
}

func BenchmarkUUIDFast(b *testing.B) {
	opt := Option{}
	opt.IndexBits = 18
	opt.WorkerIdBits = 10
	opt.TimeBits = 35
	opt.TimeLatest = false
	opt.WorkerId = 1
	opt.Epoch = internal.UnixMs() + (3600 * 24 * 397 * 1000)
	var gen IGenerator = NewFastGenerator(opt)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gen.UUID()
	}
}

// go test -bench="." -parallel 50 -count 5
func BenchmarkUUIDParallel(b *testing.B) {
	opt := Option{}
	opt.IndexBits = 18
	opt.WorkerIdBits = 10
	opt.TimeBits = 35
	opt.TimeLatest = true
	opt.WorkerId = 1
	opt.Epoch = internal.UnixMs() + (3600 * 24 * 397 * 1000)
	var gen IGenerator = NewGenerator(opt)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) { //并发
		for pb.Next() {
			gen.UUID()
			//b.Fail()
		}
	})
}

// go test -bench="." -parallel 50 -count 5
// go test -bench="." -parallel 500000 -count 1 -benchmem -cpu 2,4,8,16
func BenchmarkUUIDFastParallel(b *testing.B) {
	opt := Option{}
	opt.IndexBits = 18
	opt.WorkerIdBits = 10
	opt.TimeBits = 35
	opt.TimeLatest = true
	opt.WorkerId = 1
	opt.Epoch = internal.UnixMs() + (3600 * 24 * 397 * 1000)
	var gen IGenerator = NewFastGenerator(opt)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) { //并发
		for pb.Next() {
			gen.UUID()
		}
	})
}
