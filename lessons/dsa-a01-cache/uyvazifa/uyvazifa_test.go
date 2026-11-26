// go test -run '^$' -bench . -count 10 ./uyvazifa   (taskset -c 2 bilan)
//
// Diqqat: o'lchanadigan siklni to'g'ridan-to'g'ri `for b.Loop() { ... }` ichiga yozmang.
// b.Loop tanasidagi o'zgaruvchilarni "tirik" ushlaydi (go doc testing.B.Loop), shunda `sum`
// registrda emas, stack'da yashaydi va sikl bir necha barobar sekinlashadi. Funksiyaga chiqaring.
package uyvazifa

import (
	"fmt"
	"testing"
)

const n = 1 << 24 // 16 777 216 element: hech bir keshga sig'maydi

// toldir: har elementga yozamiz. Hech yozilmagan xotirani Linux bitta umumiy nol sahifaga
// ulaydi va uni o'qish keshdan keladi — o'lchov yolg'on bo'lardi.
func toldir[T int32 | int64](s []T) []T {
	for i := range s {
		s[i] = T(i)
	}
	return s
}

func TestYigindi(t *testing.T) {
	s64, s32 := toldir(make([]int64, 1000)), toldir(make([]int32, 1000))
	nq := make([]Nuqta, 1000)
	for i := range nq {
		nq[i].X = int64(i)
	}
	if Yig64(s64) != 499500 || Yig32(s32) != 499500 || YigX(nq) != 499500 || Qadam(s64, 1) != 499500 {
		t.Fatal("yig'indi noto'g'ri")
	}
}

// O'rta: qadam k = 1 … 32. 64 baytlik qatorda 8 ta int64 bor.
func BenchmarkQadam(b *testing.B) {
	s := toldir(make([]int64, n))
	for _, k := range []int{1, 2, 4, 8, 16, 32} {
		b.Run(fmt.Sprintf("k=%d", k), func(b *testing.B) {
			for b.Loop() {
				Qadam(s, k)
			}
		})
	}
}

// Qiyin: int32 ikki barobar kam bayt.
func BenchmarkInt64(b *testing.B) {
	s := toldir(make([]int64, n))
	for b.Loop() {
		Yig64(s)
	}
}

func BenchmarkInt32(b *testing.B) {
	s := toldir(make([]int32, n))
	for b.Loop() {
		Yig32(s)
	}
}

// Guru: struct'lar slice'i (AoS) va alohida slice (SoA). Faqat X kerak.
func BenchmarkAoS(b *testing.B) {
	s := make([]Nuqta, n)
	for i := range s {
		s[i] = Nuqta{int64(i), 1, 2, 3}
	}
	for b.Loop() {
		YigX(s)
	}
}

func BenchmarkSoA(b *testing.B) {
	xs := toldir(make([]int64, n)) // Y, Z, W boshqa slice'larda
	for b.Loop() {
		Yig64(xs)
	}
}
