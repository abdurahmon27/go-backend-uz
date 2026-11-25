// xato: ATAYLAB kompilyatsiya boʻlmaydi.
// Buyurtma emas, faqat *Buyurtma fmt.Stringer'ni bajaradi.
// go build ./xato
package main

import "fmt"

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

func (b *Buyurtma) String() string { return fmt.Sprint("buyurtma: ", b.Summa) }

func main() {
	b := Buyurtma{Summa: 150000}
	var s fmt.Stringer = b
	fmt.Println(s)
}
