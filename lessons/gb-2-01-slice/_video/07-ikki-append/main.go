package main

import "fmt"

func main() {
	stock := []int{40, 12, 7}
	stock = append(stock, 25) // len 4, cap 6
	a := append(stock, 3)
	b := append(stock, 9)
	fmt.Println(a[4], b[4])
}
