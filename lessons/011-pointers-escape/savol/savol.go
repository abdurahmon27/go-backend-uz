// savol: ikki chaqiruvchi, ikkalasi lokal buyurtmaning manzilini beradi.
// Qaysi buyurtma heap'ga chiqadi? (go:noinline: kompilyator chaqiruvchini
// faqat funksiyaning qisqa xulosasi — "parameter tags" — boʻyicha tahlil qilsin)
// go build -gcflags=-m ./savol 2>&1 | grep moved
package savol

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

//go:noinline
func summa(p *Buyurtma) int64 { return p.Summa } // p does not escape

var oxirgi *Buyurtma

//go:noinline
func eslab(p *Buyurtma) { oxirgi = p } // leaking param: p

func birinchi() int64 {
	a := Buyurtma{Summa: 100}
	return summa(&a)
}

func ikkinchi() {
	b := Buyurtma{Summa: 200}
	eslab(&b)
}
