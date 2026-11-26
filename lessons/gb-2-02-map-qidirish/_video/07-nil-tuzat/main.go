package main

import "fmt"

func main() {
	stock := make(map[string]int)
	fmt.Println(stock["olma"], len(stock))
	stock["olma"] = 40
	fmt.Println(stock["olma"])
}
