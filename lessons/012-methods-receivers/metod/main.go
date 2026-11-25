// metod: metod — receiver'li funksiya. Value receiver nusxa oladi, pointer receiver manzil oladi.
package main

import "fmt"

type Buyurtma struct {
	Summa   int64 // tiyinda
	Tolandi bool
}

// Jami: value receiver. b — buyurtmaning nusxasi.
func (b Buyurtma) Jami() int64 { return b.Summa + b.Summa/100 } // +1% xizmat haqi

// Tola: value receiver. Nusxani oʻzgartiradi, asl buyurtma oʻzgarmaydi.
func (b Buyurtma) Tola() { b.Tolandi = true }

// TolaP: pointer receiver. Asl buyurtmani oʻzgartiradi.
func (b *Buyurtma) TolaP() { b.Tolandi = true }

func main() {
	b := Buyurtma{Summa: 50_000}
	fmt.Println("Jami:          ", b.Jami())
	fmt.Println("Buyurtma.Jami: ", Buyurtma.Jami(b)) // receiver — birinchi argument

	b.Tola()
	fmt.Println("Tola:  ", b.Tolandi)

	b.TolaP() // (&b).TolaP() ning qisqa yozuvi
	fmt.Println("TolaP: ", b.Tolandi)
}
