// NaN: sharhi yoʻq mahsulotning reytingi 0/0 = NaN. Qoʻlda yozilgan taqqoslash buziladi.
package main

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

type Mahsulot struct {
	Nom     string
	Reyting float64
}

func main() {
	var nol float64
	yangi := Mahsulot{"yangi", nol / nol} // 0/0 = NaN
	asl := []Mahsulot{
		{"dazmol", 4.8}, {"fen", 3.9}, yangi,
		{"choynak", 4.2}, {"mikser", 4.5}, {"blender", 4.5},
	}

	// Xato: NaN hech narsadan katta ham, kichik ham emas — "hammaga teng".
	mm := slices.Clone(asl)
	slices.SortFunc(mm, func(a, b Mahsulot) int {
		if a.Reyting > b.Reyting {
			return -1
		}
		if a.Reyting < b.Reyting {
			return 1
		}
		return 0
	})
	fmt.Println("qoʻlda:", nomlar(mm))

	// Toʻgʻri: cmp.Compare NaN'ni eng kichik deb hisoblaydi.
	// cmp.Or: reyting teng boʻlsa, nom boʻyicha.
	mm = slices.Clone(asl)
	slices.SortFunc(mm, func(a, b Mahsulot) int {
		return cmp.Or(
			cmp.Compare(b.Reyting, a.Reyting),
			strings.Compare(a.Nom, b.Nom),
		)
	})
	fmt.Println("cmp.Or:", nomlar(mm))
}

func nomlar(mm []Mahsulot) []string {
	var s []string
	for _, m := range mm {
		s = append(s, fmt.Sprintf("%s=%.1f", m.Nom, m.Reyting))
	}
	return s
}
