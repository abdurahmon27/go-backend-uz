package main

import "fmt"

func main() {
	var x int8 = 127
	fmt.Printf("%4d = %08b\n", x, uint8(x))
	x++
	fmt.Printf("%4d = %08b\n", x, uint8(x))
	x = -1
	fmt.Printf("%4d = %08b\n", x, uint8(x))
}
