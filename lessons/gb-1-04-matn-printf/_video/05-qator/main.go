package main

import "fmt"

func main() {
	name := "choy"
	var price int64 = 199_900 // tiyinda
	fmt.Printf("%s: %d tiyin", name, price)
	fmt.Printf("%s: %d tiyin", "non", 400_000)
}
