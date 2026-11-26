package main

import "fmt"

func main() {
	var stock map[string]int
	fmt.Println(stock["olma"], len(stock))
	stock["olma"] = 40
	fmt.Println(stock["olma"])
}
