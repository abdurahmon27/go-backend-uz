// hook: funksiya nil qaytaradi, ekranga <nil> chiqadi, lekin err != nil.
// ATAYLAB XATO (typed nil). Tuzatilgani: ../fix
// go run ./hook
package main

import "fmt"

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

type SummaXato struct{ Summa int64 }

func (e *SummaXato) Error() string { return fmt.Sprintf("summa notoʻgʻri: %d", e.Summa) }

func tekshir(b Buyurtma) error {
	var xato *SummaXato
	if b.Summa <= 0 {
		xato = &SummaXato{Summa: b.Summa}
	}
	return xato
}

func main() {
	err := tekshir(Buyurtma{Summa: 150000})
	fmt.Println("err:", err)
	if err != nil {
		fmt.Println("xato bor!")
	}
	fmt.Printf("ichida: %T, %v\n", err, err)
}
