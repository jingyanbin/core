package basal

import "testing"

func BenchmarkFormatString(b *testing.B) {
	s := "this is a test @<0>"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Format(s, '@', '<', '>', KwArgs{"0": 99}, false)
	}
}
