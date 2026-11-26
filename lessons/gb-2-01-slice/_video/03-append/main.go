package main

import "fmt"

func main() {
	stock := []int{40, 12, 7}
	fmt.Printf("%d %d %p\n", len(stock), cap(stock), stock)
	stock = append(stock, 25) // oʻrik
	fmt.Printf("%d %d %p\n", len(stock), cap(stock), stock)
	stock = append(stock, 3) // anor
	fmt.Printf("%d %d %p\n", len(stock), cap(stock), stock)
}
