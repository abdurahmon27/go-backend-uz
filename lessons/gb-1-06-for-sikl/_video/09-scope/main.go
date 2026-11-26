package main

import "fmt"

func main() {
	sold := 3
	var subtotal int64
	for i := 1; i <= sold; i++ {
		subtotal += 199_900
	}
	fmt.Println(subtotal, i)
}
