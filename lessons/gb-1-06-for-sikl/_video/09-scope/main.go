package main

import "fmt"

func main() {
	qty := 3
	var subtotal int64
	for i := 1; i <= qty; i++ {
		subtotal += 199_900
	}
	fmt.Println(subtotal, i)
}
