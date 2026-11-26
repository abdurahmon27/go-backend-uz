package main

import "fmt"

func main() {
	stock := []int{40, 12, 7} // olma, non, choy
	fmt.Println(stock, len(stock), cap(stock))
	fmt.Println(stock[0], stock[2])
}
