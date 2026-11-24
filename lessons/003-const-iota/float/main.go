package main

import "fmt"

const a, b = 0.1, 0.2

func main() {
	fmt.Println(a+b == 0.3)

	x, y := 0.1, 0.2
	fmt.Println(x+y == 0.3)
	fmt.Println(x + y)
}
