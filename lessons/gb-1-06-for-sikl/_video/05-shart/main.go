package main

import "fmt"

func main() {
	var price int64 = 199_900    // tiyinda
	var budget int64 = 5_000_000 // 50 000 soʻm, tiyinda

	var spent int64
	count := 0
	for spent+price <= budget {
		spent += price
		count++
	}
	fmt.Println(count, "dona")
}
