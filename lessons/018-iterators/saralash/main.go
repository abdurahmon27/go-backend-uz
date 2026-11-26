// Saralash: reyting boʻyicha kamayish tartibida. Ikki xato va bitta toʻgʻri usul.
package main

import (
	"cmp"
	"fmt"
	"slices"
)

type Mahsulot struct {
	Nom     string
	Reyting float64
}

func main() {
	mm := []Mahsulot{{"choynak", 4.2}, {"dazmol", 4.8}, {"fen", 3.9}, {"blender", 4.5}}

	// Xato: ayirma int'ga aylanganda 0.6 → 0 boʻladi, ya'ni "teng".
	slices.SortFunc(mm, func(a, b Mahsulot) int {
		return int(b.Reyting - a.Reyting)
	})
	fmt.Println("int(ayirma):  ", mm)

	// Toʻgʻri: cmp.Compare manfiy, nol yoki musbat qaytaradi.
	slices.SortFunc(mm, func(a, b Mahsulot) int {
		return cmp.Compare(b.Reyting, a.Reyting)
	})
	fmt.Println("cmp.Compare:  ", mm)
}
