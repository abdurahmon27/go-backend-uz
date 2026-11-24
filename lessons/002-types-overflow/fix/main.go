package main

import "fmt"

func main() {
	var stock uint = 3
	order := uint(5)

	if order > stock {
		fmt.Println("yetarli emas: bor", stock, "ta, so'ralgan", order, "ta")
		return
	}
	stock -= order
	fmt.Println("qoldiq:", stock)
}
