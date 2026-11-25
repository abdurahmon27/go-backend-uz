// Uy vazifasi (oson): Mahsulot turi fmt.Stringer'ni bajarsin, buni kompilyator tekshirsin.
// go run ./homework/oson
package main

import "fmt"

type Mahsulot struct {
	Nom  string
	Narx int64
}

func (m Mahsulot) String() string { return fmt.Sprintf("%s (%d soʻm)", m.Nom, m.Narx) }

// Kompilyatsiya vaqtidagi tekshiruv (Effective Go: Interface checks).
// String metodini oʻchirsangiz, dastur kompilyatsiya boʻlmaydi.
var _ fmt.Stringer = Mahsulot{}

func main() {
	fmt.Println(Mahsulot{Nom: "choynak", Narx: 85000})
}
