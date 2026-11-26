package main

import "fmt"

func main() {
	prices := []int{2, 3, 5}
	total := 0
	for p := range prices {
		total += p
	}
	fmt.Println(total)
}
