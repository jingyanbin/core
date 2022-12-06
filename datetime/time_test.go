package datetime

import (
	"testing"
)
import "time"

//go test -bench="." -count 5
//func BenchmarkUnix1(b *testing.B) {
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		Unix()
//	}
//}
//
//func BenchmarkUnix2(b *testing.B) {
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		time.Now().Unix()
//	}
//}

//func BenchmarkDateTime(b *testing.B) {
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		dt := Now()
//		dt.Year()
//		dt.Month()
//		dt.Day()
//		dt.Min()
//		dt.Sec()
//	}
//}
//
//func BenchmarkTime(b *testing.B) {
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		dt := time.Now()
//		dt.Year()
//		dt.Month()
//		dt.Day()
//		dt.Hour()
//		dt.Minute()
//		dt.Second()
//	}
//}

////150 ns/op
//func BenchmarkDateTimeFormat(b *testing.B) {
//	dt := Now()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		dt.YmdHMS()
//	}
//}
//
////190 ns/op
//func BenchmarkTimeFormat(b *testing.B) {
//	dt := time.Now()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		dt.Format("2006-01-02 15:04:05")
//	}
//}

// 90 ns/op
func BenchmarkDateTimeFormatToUnix(b *testing.B) {
	dtStr := "2021-11-23 16:30:00"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		YmdHMSToUnix(dtStr, Local(), true)
	}
}

// 175 ns/op
func BenchmarkTimeFormatToUnix(b *testing.B) {
	dtStr := "2021-11-23 16:30:00"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		time.ParseInLocation("2006-01-02 15:04:05", dtStr, time.Local)
	}
}
