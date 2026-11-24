package kabisa

import "testing"

// natija paket darajasidagi oʻzgaruvchiga yoziladi: kompilyator
// hisobni oʻchira olmaydi, chunki natija funksiyadan tashqarida koʻrinadi.
var natija bool

func BenchmarkKabisaBN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		kabisa(i)
	}
}

func BenchmarkSinkBN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		natija = kabisa(i)
	}
}

func BenchmarkKabisaLoop(b *testing.B) {
	yil := 0
	for b.Loop() {
		kabisa(yil)
		yil++
	}
}
