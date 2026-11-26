package main

import "fmt"

func main() {
	var price int64 = 199_900 // tiyinda
	qty := 3

	var subtotal int64
	for i := 1; i <= qty; i++ {
		subtotal += price
	}
	fmt.Println(subtotal)
}
