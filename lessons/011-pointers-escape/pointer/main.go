// pointer: qiymat nusxa boʻlib uzatiladi, pointer esa asl qiymatga yoʻl beradi.
package main

import "fmt"

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

// tola: b — nusxa. Oʻzgarish chaqiruvchiga yetmaydi.
func tola(b Buyurtma) {
	b.Tolandi = true
}

// tolaP: p — manzil. *p orqali asl buyurtmani oʻzgartiramiz.
func tolaP(p *Buyurtma) {
	p.Tolandi = true // (*p).Tolandi ning qisqa yozuvi
}

func main() {
	b := Buyurtma{Summa: 50_000}

	tola(b)
	fmt.Println("tola:  ", b.Tolandi)

	tolaP(&b)
	fmt.Println("tolaP: ", b.Tolandi)

	var p *Buyurtma // zero value: nil
	fmt.Println("p == nil:", p == nil)
}
