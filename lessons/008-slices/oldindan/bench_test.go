package oldindan

import "testing"

func BenchmarkNil(b *testing.B) {
	for b.Loop() {
		yigNil(1000)
	}
}

func BenchmarkMake(b *testing.B) {
	for b.Loop() {
		yigMake(1000)
	}
}
