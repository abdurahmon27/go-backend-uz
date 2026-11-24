package kabisa

import "testing"

func BenchmarkBoshBN(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}

func BenchmarkKabisaBN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		kabisa(i)
	}
}
