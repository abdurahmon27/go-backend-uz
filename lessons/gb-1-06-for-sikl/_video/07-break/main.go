package main

import "fmt"

func main() {
	var price int64 = 199_900 // tiyinda
	qty := 7

	var subtotal int64
	for i := 1; i <= qty; i++ {
		if i%3 == 0 {
			break
		}
		subtotal += price
	}
	fmt.Println(subtotal)
}
