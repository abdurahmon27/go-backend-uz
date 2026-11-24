package kabisa

import "testing"

func BenchmarkBoshLoop(b *testing.B) {
	yil := 0
	for b.Loop() {
		yil++
	}
}

func BenchmarkKabisaLoop(b *testing.B) {
	yil := 0
	for b.Loop() {
		kabisa(yil)
		yil++
	}
}
