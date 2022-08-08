package deepcopy

//go test -bench=Bench* -v .\copy_test.go .\copy.go .\copy_other.go -benchmem -count=3
import (
	"testing"
)

type Test struct {
	A1 int64
	A2 float64
	A3 bool
	A4 string
	A5 []byte
}

func NewTest() *Test {
	t5 := &Test{}
	t5.A1 = 1    //proto.Int64(1)   //1
	t5.A2 = 2    //proto.Int64(2)   //2
	t5.A3 = true //proto.Float32(3) //3
	t5.A4 = "xx" //proto.Float32(4) // 4
	return t5
	//
	//t5.A5 = append(t5.A5, 1, 2, 3, 4, 5, 6, 7, 8, 9, 7, 5, 4, 5, 5, 5, 5, 6, 5, 2, 3, 2, 1, 2, 31, 21, 21, 21, 1, 2, 1)
	//t6 := t5
	//t7 := t5
	//
	//t4 := &pb.Test2{}
	//t4.A1 = 1 //proto.Int64(1)   //1
	//t4.A2 = 2 //proto.Int64(1)   //2
	//t4.A3 = 3 //proto.Float32(1) //3
	//t4.A4 = &t6
	//t4.A5 = append(t4.A5, &t5, &t5, &t5, &t5, &t5, &t5, &t5, &t5, &t5, &t5)
	//t4.A6 = &t7
	//
	//t3 := &pb.Test{}
	//t3.A1 = 1 //proto.Int64(1)   //1
	//t3.A2 = 2 //proto.Int64(1)   //2
	//t3.A3 = 3 //proto.Float32(1) //3
	//t3.A4 = 4 //proto.Float32(1) //4
	//t3.A5 = append(t3.A5, t4)
	//t3.A6 = append(t3.A5, t4, t4, t4, t4, t4, t4)
	//t3.A7 = t4

	//t2 := &pb.ForwardMessage{}
	//t2.Ip = 111
	//t2.ClientId = 1
	//t2.Content = []byte("sadasda")
	//t2.UserId = 5142132
	//t2.ServerId = 359884
	//return nil
	//t1 := &Test{}
	//t1.A1 = 1
	//t1.A2 = 2
	//t1.A3 = true
	//t1.A4 = "111"
	//t1.A5 = []byte{1, 2, 3}
	//return t1
	return nil
}

func BenchmarkCopyByJson(b *testing.B) {
	t1 := NewTest()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CopyByJson(t1)
	}
}

func BenchmarkCopyByMsgPack(b *testing.B) {
	t1 := NewTest()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CopyByMsgPack(t1)
	}
}

func BenchmarkCopyByGob(b *testing.B) {
	t1 := NewTest()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CopyByGob(t1)

	}
}

func BenchmarkCopy(b *testing.B) {
	t1 := NewTest()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Copy(t1)
	}
}

//func BenchmarkCopyByGoGo(b *testing.B) {
//	t1 := NewTest()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		CopyByGoGo(t1)
//	}
//}
