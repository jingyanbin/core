package uuid

import (
	"math/rand"
	"testing"
)

func TestUUID(t *testing.T) {
	//gen := NewHexGenerator(1002, 0, 0, true)
	gen := NewUUIDGenerator(65536, 16, 0, true, false)
	//gen1 := NewHexGenerator(2, 0, 0, true)
	t.Logf("===%d", rand.Intn(1))
	//var index int

	for i:=0; i < 10; i++{
		uid := gen.UUID()
		ms, index, serverId := gen.DeUUID(uid)
		t.Logf("ms: %v, index: %v, serverId: %v", ms, index, serverId)


		//ms, index, workerId, err := gen.DeUUID(uid)
		//
		//t.Logf("uuid: %s, ms:%d,index：%d, worker: %d, err: %v", uid, ms, index, workerId, err)
		//uid2 := gen.UUIDExtra(uint16(i))
		//ms2, index2, workerId2, extra, err := gen.DeUUIDExtra(uid2)

		//t.Logf("uuid:2 %s, ms:%d,index：%d, worker: %d, extra: %d, err: %v", uid2, ms2, index2, workerId2, extra, err)

		//ms, index, workerId, extra, err := gen1.DeUUIDExtra(uid)

	}
}

func BenchmarkUUID(b *testing.B) {
	//var a int64 = -1
	//64*1000
	//gen := NewHexGenerator(1, 0, 0, true)
	gen := NewUUIDGenerator(1, 16, 0, true, false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gen.UUID()
	}
}

//go test -bench="." -parallel 50 -count 5
func BenchmarkUUIDParallel(b *testing.B) {

	gen := NewHexGenerator(1, 0, 0, true)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) { //并发
		for pb.Next() {
			gen.UUID()

			//b.Fail()
		}
	})
}

//go test -bench="." -parallel 50 -count 5
func BenchmarkInt64UUIDParallel(b *testing.B) {
	gen := NewUUIDGenerator(101, 9, -1, true, false)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) { //并发
		for pb.Next() {
			gen.UUID()
			//b.Fail()
		}
	})
}