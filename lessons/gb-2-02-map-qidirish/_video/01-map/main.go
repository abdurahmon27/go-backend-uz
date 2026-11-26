package main

import "fmt"

const tiyinPerSom = 100 // 1 soʻm = 100 tiyin

func main() {
	prices := map[string]int64{
		"olma": 12_000 * tiyinPerSom,
		"non":  4_000 * tiyinPerSom,
	}
	fmt.Println(prices["olma"])
	prices["choy"] = 25_000 * tiyinPerSom
	fmt.Println(len(prices))
	fmt.Println(prices["anor"])
}
