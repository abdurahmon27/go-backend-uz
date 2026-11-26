package main

import "fmt"

func main() {
	prices := []int{2, 3, 5}
	total := 0
	for _, p := range prices { // TS: for (const p of prices)
		total += p
	}
	fmt.Println(total) // 10
}
