package main

import "fmt"

func main() {
	narxlar := []int{100, 200, 300}
	for _, n := range narxlar {
		n *= 2
	}
	fmt.Println(narxlar)

	for i := range narxlar {
		narxlar[i] *= 2
	}
	fmt.Println(narxlar)
}
