package main

import "fmt"

func main() {
	var price int64 = 199_900 // tiyinda
	sold := 3

	var subtotal int64
	for range sold {
		subtotal += price
	}
	fmt.Println(subtotal)
}
