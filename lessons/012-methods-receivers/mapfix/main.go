// mapfix: map elementida pointer metodni chaqirishning ikki yoʻli.
package main

import "fmt"

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

func (b *Buyurtma) Tola() { b.Tolandi = true }

func main() {
	// 1-yoʻl: nusxa ol, oʻzgartir, qaytarib yoz.
	m := map[string]Buyurtma{"a1": {Summa: 50_000}}
	b := m["a1"]
	b.Tola()
	m["a1"] = b
	fmt.Println("1-yoʻl:", m["a1"].Tolandi)

	// 2-yoʻl: map'da pointer saqla. m2["a1"] — pointer, uning ortidagi buyurtma addressable.
	m2 := map[string]*Buyurtma{"a1": {Summa: 50_000}}
	m2["a1"].Tola()
	fmt.Println("2-yoʻl:", m2["a1"].Tolandi)
}
