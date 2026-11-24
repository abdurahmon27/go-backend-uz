package main

import "fmt"

func main() {
	var x int16 = 32767 // math.MaxInt16
	x++
	fmt.Println(x) // -32768: 0111111111111111 + 1 = 1000000000000000
}
