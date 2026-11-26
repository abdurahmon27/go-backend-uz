// Hook: Tolanganlar slice qaytarmaydi, u funksiya qaytaradi.
// for range esa shu funksiya ustida aylanadi (Go 1.23+).
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

// iter paketida: type Seq[V any] func(yield func(V) bool)
//
// Tolanganlar faqat toʻlangan buyurtmalarni beradi.
func Tolanganlar(bb []Buyurtma) iter.Seq[Buyurtma] {
	return func(yield func(Buyurtma) bool) {
		for _, b := range bb {
			if !b.Tolandi {
				continue
			}
			if !yield(b) {
				return
			}
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
