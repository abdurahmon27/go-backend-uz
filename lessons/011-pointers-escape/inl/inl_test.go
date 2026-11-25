package inl

import "testing"

func BenchmarkQiymat(b *testing.B) {
	for b.Loop() {
		_ = yarat(100)
	}
}

func BenchmarkPointer(b *testing.B) {
	for b.Loop() {
		_ = yaratP(100)
	}
}
