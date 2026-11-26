package kesh

import (
	"math/rand/v2"
	"testing"
)

func yangiRand() *rand.Rand { return rand.New(rand.NewPCG(1, 2)) }

func TestBirXilYigindi(t *testing.T) {
	for _, n := range []int{1, 2, 7, 1 << 10, 1 << 16} {
		want := int64(n) * int64(n-1) / 2
		s := Sonlar(n)
		got := []int64{Ketma(s), Aralash(s, Tartib(n, yangiRand())), Zanjir(ZanjirSlice(n, yangiRand())), Zanjir(TartibliZanjir(n)),
			Royxat(RoyxatYangi(n)), Royxat(RoyxatSochilgan(n, yangiRand()))}
		for i, g := range got {
			if g != want {
				t.Fatalf("n=%d, %d-usul: %d, kutilgan %d", n, i, g, want)
			}
		}
	}
	if got := Qadam(Sonlar(16), 8); got != 0+8 {
		t.Fatalf("Qadam: %d", got)
	}
}

func BenchmarkKetma(b *testing.B) {
	olcha(b, func(n int) func() int64 {
		s := Sonlar(n)
		return func() int64 { return Ketma(s) }
	})
}

func BenchmarkAralash(b *testing.B) {
	olcha(b, func(n int) func() int64 {
		s, t := Sonlar(n), Tartib(n, yangiRand())
		return func() int64 { return Aralash(s, t) }
	})
}

func BenchmarkTZanjir(b *testing.B) {
	olcha(b, func(n int) func() int64 {
		s := TartibliZanjir(n)
		return func() int64 { return Zanjir(s) }
	})
}

func BenchmarkZanjir(b *testing.B) {
	olcha(b, func(n int) func() int64 {
		s := ZanjirSlice(n, yangiRand())
		return func() int64 { return Zanjir(s) }
	})
}
