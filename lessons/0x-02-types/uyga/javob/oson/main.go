// Javob 🟢. Uchta xato edi:
//  1. total := 0 → int; float64 qo'shib bo'lmaydi → total := 0.0 (yoki var total float64)
//  2. p*qty: float64 * int → float64(qty)
//  3. i ishlatilmagan (declared and not used) → _
//
// Natija: jami: 40
package main

import "fmt"

func main() {
	prices := []float64{12.5, 7.5}
	qty := 2
	total := 0.0
	for _, p := range prices {
		total = total + p*float64(qty)
	}
	fmt.Println("jami:", total)
}
