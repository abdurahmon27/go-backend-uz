// Package guru: uy vazifasi (guru) javobi. Qaysi instantiatsiyalar bitta kompilyatsiya
// qilingan tanani boʻlishadi? go tool nm bilan tekshiring (README ga qarang).
package guru

type Summali interface{ Jami() int64 }

type Buyurtma struct{ Summa int64 }

func (b *Buyurtma) Jami() int64 { return b.Summa }

type Mijoz struct{ Qarz int64 }

func (m *Mijoz) Jami() int64 { return m.Qarz }

// Narx — pointer emas, asosiy turi int64.
type Narx int64

func (n Narx) Jami() int64 { return int64(n) }

// Chegirma — u ham int64 asosli, lekin boshqa tur.
type Chegirma int64

func (c Chegirma) Jami() int64 { return -int64(c) }

//go:noinline
func JamiGeneric[T Summali](xs []T) (s int64) {
	for _, x := range xs {
		s += x.Jami()
	}
	return s
}
