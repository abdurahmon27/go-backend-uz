// assert: type assertion — interface ichidagi turni tekshirish.
// go run ./assert
package main

import "fmt"

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

func main() {
	var a any = Buyurtma{Summa: 150000}

	b, ok := a.(Buyurtma) // comma-ok: panic yoʻq
	fmt.Println("Buyurtma?", ok, b.Summa)

	s, ok := a.(string)
	fmt.Printf("string?   %v %q\n", ok, s)

	n := a.(int) // ok'siz: tur mos kelmasa — panic
	fmt.Println(n)
}
