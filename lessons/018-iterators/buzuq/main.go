// Buzuq: hook/main.go bilan bir xil, faqat iterator yield natijasini tekshirmaydi.
package main

import (
	"fmt"
	"iter"
)

type Buyurtma struct {
	ID      int
	Summa   int64
	Tolandi bool
}

// Tolanganlar — XATO: yield false qaytarsa ham davom etadi.
func Tolanganlar(bb []Buyurtma) iter.Seq[Buyurtma] {
	return func(yield func(Buyurtma) bool) {
		for _, b := range bb {
			if !b.Tolandi {
				continue
			}
			yield(b)
		}
	}
}

func main() {
	bb := []Buyurtma{
		{1, 150000, true},
		{2, 99000, false},
		{3, 42000, true},
		{4, 310000, true},
	}
	for b := range Tolanganlar(bb) {
		fmt.Println("buyurtma", b.ID, b.Summa)
		if b.Summa < 50000 {
			fmt.Println("kichik summa, toʻxtaymiz")
			break
		}
	}
}
