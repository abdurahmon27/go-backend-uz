// juftlik: interface qiymati = (tur, qiymat) juftligi.
// go run ./juftlik
package main

import (
	"fmt"
	"unsafe"
)

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

func korsat(nom string, a any) {
	qiymat := fmt.Sprint(a)
	fmt.Printf("%-9s tur: %-15T qiymat: %-14s a == nil: %v\n", nom, a, qiymat, a == nil)
}

func main() {
	var p *Buyurtma // nil pointer

	korsat("nil", nil)
	korsat("42", 42)
	korsat("Buyurtma", Buyurtma{Summa: 150000})
	korsat("p", p)

	var a any
	fmt.Println("any hajmi:", unsafe.Sizeof(a), "bayt")
}
