package main

import "fmt"

func main() {
	stock := []int{40, 12, 7}
	append(stock, 25)
	fmt.Println(stock)
}
