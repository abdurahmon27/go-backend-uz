package main

import "fmt"

const tiyinPerSom = 100 // 1 soʻm = 100 tiyin

func main() {
	prices := map[string]int64{
		"olma": 12_000 * tiyinPerSom,
		"non":  4_000 * tiyinPerSom,
	}
	prices["olma"] = 0
	_, ok := prices["olma"]
	fmt.Println(len(prices), ok)
	delete(prices, "non")
	_, ok = prices["non"]
	fmt.Println(len(prices), ok)
}
