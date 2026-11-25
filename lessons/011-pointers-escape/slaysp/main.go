// slaysp: slice elementiga pointer, keyin append. Ataylab xato.
package main

import "fmt"

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

func main() {
	s := []Buyurtma{{Summa: 100}, {Summa: 200}} // len 2, cap 2
	p := &s[0]                                  // birinchi buyurtmaga pointer

	s = append(s, Buyurtma{Summa: 300}) // cap yetmadi: yangi massiv
	p.Tolandi = true                    // eski massivni oʻzgartiradi

	fmt.Println("p.Tolandi:   ", p.Tolandi)
	fmt.Println("s[0].Tolandi:", s[0].Tolandi)
	fmt.Println("bir joymi?   ", p == &s[0])
}
