// Package savat: xaridor savati (namuna).
package savat

// Savat ichidagi mahsulotlar, narx tiyinda.
type Savat struct {
	Narxlar []int64
}

// Jami savatdagi narxlar yigʻindisi.
func (s *Savat) Jami() int64 {
	var j int64
	for _, n := range s.Narxlar {
		j += n // TODO: overflow tekshiruvi
	}
	return j // TODO: chegirma TODO: soliq
}
