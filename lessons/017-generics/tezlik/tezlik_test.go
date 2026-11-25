package tezlik

import "testing"

const n = 1000

func buyurtmalar() []*Buyurtma {
	xs := make([]*Buyurtma, n)
	for i := range xs {
		xs[i] = &Buyurtma{Summa: int64(i)}
	}
	return xs
}

func TestBirXil(t *testing.T) {
	xs := buyurtmalar()
	is := make([]Summali, n)
	for i, x := range xs {
		is[i] = x
	}
	a, b, c := JamiInterface(is), JamiGeneric(xs), JamiOddiy(xs)
	if a != b || b != c {
		t.Fatalf("%d %d %d", a, b, c)
	}
}

func BenchmarkInterface(b *testing.B) {
	xs := buyurtmalar()
	is := make([]Summali, n)
	for i, x := range xs {
		is[i] = x
	}
	for b.Loop() {
		JamiInterface(is)
	}
}

func BenchmarkGeneric(b *testing.B) {
	xs := buyurtmalar()
	for b.Loop() {
		JamiGeneric(xs)
	}
}

func BenchmarkOddiy(b *testing.B) {
	xs := buyurtmalar()
	for b.Loop() {
		JamiOddiy(xs)
	}
}
