package kabisa

import "testing"

// Tajriba (README): doimiy argument. Go 1.26 kabisa'ni inline qilib,
// kabisa(2024) = true javobini kompilyatsiya paytida hisoblaydi.
// b.Loop natijani "tirik" saqlaydi, lekin hisobning oʻzi siklda qolmaydi.
func BenchmarkDoimiyLoop(b *testing.B) {
	for b.Loop() {
		kabisa(2024)
	}
}

// Har aylanishda boshqa yil: hisob haqiqatan bajariladi.
func BenchmarkKabisaLoop(b *testing.B) {
	yil := 0
	for b.Loop() {
		kabisa(yil)
		yil++
	}
}
