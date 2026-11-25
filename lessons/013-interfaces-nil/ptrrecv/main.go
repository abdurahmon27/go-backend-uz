// ptrrecv: String pointer receiver bilan. Buyurtma qiymatining method set'ida
// String yoʻq (GO-1.12), shuning uchun fmt uni koʻrmaydi.
// go run ./ptrrecv
package main

import "fmt"

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

func (b *Buyurtma) String() string {
	return fmt.Sprintf("buyurtma: %d soʻm, toʻlandi: %v", b.Summa, b.Tolandi)
}

func main() {
	b := Buyurtma{Summa: 150000}
	fmt.Println(b)  // qiymat: String ishlamaydi
	fmt.Println(&b) // pointer: String ishlaydi
}
