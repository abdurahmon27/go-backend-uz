package main

import "fmt"

func main() {
	stock := []int{40, 12, 7}
	before := stock
	stock = append(stock, 25)
	stock[0] = 38
	fmt.Println(before[0], stock[0])
}
