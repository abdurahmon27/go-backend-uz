package main

import "fmt"

func main() {
	stock := []int{40, 12, 7}
	before := stock
	stock[0] = 38
	fmt.Println(before[0])
}
