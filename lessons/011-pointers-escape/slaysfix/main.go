// slaysfix: pointer'ni append'dan keyin oling yoki indeks bilan ishlang.
package main

import "fmt"

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

func main() {
	s := []Buyurtma{{Summa: 100}, {Summa: 200}}

	s = append(s, Buyurtma{Summa: 300})
	s[0].Tolandi = true // indeks orqali: har doim hozirgi massiv

	fmt.Println("s[0].Tolandi:", s[0].Tolandi)
}
