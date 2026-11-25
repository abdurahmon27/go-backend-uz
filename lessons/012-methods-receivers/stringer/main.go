// stringer: String metodi pointer receiver'da. fmt.Println(b) uni chaqirmaydi, fmt.Println(&b) chaqiradi.
package main

import "fmt"

type Buyurtma struct {
	Summa   int64 // tiyinda
	Tolandi bool
}

// String: *Buyurtma ning metodi. Buyurtma (pointer emas) ning method set'ida yoʻq.
func (b *Buyurtma) String() string {
	return fmt.Sprintf("buyurtma: %d soʻm, toʻlandi: %t", b.Summa/100, b.Tolandi)
}

func main() {
	b := Buyurtma{Summa: 50_000}
	fmt.Println(b)  // Buyurtma fmt.Stringer emas: maydonlar chiqadi
	fmt.Println(&b) // *Buyurtma fmt.Stringer: String() chaqiriladi
}
