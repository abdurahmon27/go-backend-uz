// hook: bitta metod, bitta tur. Slice elementida chaqiruv ishlaydi, map elementida esa
// dastur kompilyatsiya ham boʻlmaydi. Bu papka ATAYLAB kompilyatsiya boʻlmaydi: go build ./hook
package main

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

func (b *Buyurtma) Tola() { b.Tolandi = true }

func main() {
	s := []Buyurtma{{Summa: 50_000}}
	s[0].Tola() // ishlaydi: slice elementi addressable

	m := map[string]Buyurtma{"a1": {Summa: 50_000}}
	m["a1"].Tola() // xato: map elementining manzilini olib boʻlmaydi
}
