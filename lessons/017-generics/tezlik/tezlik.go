// Package tezlik: bir xil yigʻindi uch usulda — interface, generic, oddiy tur.
// go test -bench . ./tezlik
package tezlik

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

func (b *Buyurtma) Jami() int64 { return b.Summa }

type Summali interface{ Jami() int64 }

//go:noinline
func JamiInterface(xs []Summali) (s int64) {
	for _, x := range xs {
		s += x.Jami()
	}
	return s
}

//go:noinline
func JamiGeneric[T Summali](xs []T) (s int64) {
	for _, x := range xs {
		s += x.Jami()
	}
	return s
}

//go:noinline
func JamiOddiy(xs []*Buyurtma) (s int64) {
	for _, x := range xs {
		s += x.Jami()
	}
	return s
}
