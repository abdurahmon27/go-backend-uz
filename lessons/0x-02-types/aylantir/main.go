package main

import "fmt"

func main() {
	a := 2
	b := 1.5
	fmt.Printf("%T %T\n", a, b) // int float64
	fmt.Println(float64(a) + b) // 3.5
}
