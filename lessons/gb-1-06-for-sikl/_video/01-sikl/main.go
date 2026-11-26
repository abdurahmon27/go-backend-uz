package main

import "fmt"

func main() {
	var price int64 = 199_900 // tiyinda
	sold := 3

	var subtotal int64
	for i := 1; i <= sold; i++ {
		subtotal += price
	}
	fmt.Println(subtotal)
}
