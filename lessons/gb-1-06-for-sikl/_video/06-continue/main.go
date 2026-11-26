package main

import "fmt"

func main() {
	var price int64 = 199_900 // tiyinda
	sold := 7

	var subtotal int64
	for i := 1; i <= sold; i++ {
		if i%3 == 0 {
			continue // har uchinchi dona bepul
		}
		subtotal += price
	}
	fmt.Println(subtotal)
}
