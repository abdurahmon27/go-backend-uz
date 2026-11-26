// Uy vazifasi 🟢: savat summasi build bo'lmaydi. Uchta xatoni toping va tuzating.
// Avval natijani taxmin qiling, keyin: go run uyga/oson/main.go
package main

import "fmt"

func main() {
	prices := []float64{12.5, 7.5} // so'mda emas, misol uchun
	qty := 2
	total := 0
	for i, p := range prices {
		total = total + p*qty
	}
	fmt.Println("jami:", total)
}
