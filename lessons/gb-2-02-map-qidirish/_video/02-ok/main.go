package main

import (
	"fmt"
	"os"
)

const tiyinPerSom = 100 // 1 soʻm = 100 tiyin

func main() {
	prices := map[string]int64{
		"olma": 12_000 * tiyinPerSom,
		"non":  4_000 * tiyinPerSom,
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "foydalanish: bozor NOM")
		os.Exit(1)
	}
	name := os.Args[1]
	price, ok := prices[name]
	if !ok {
		fmt.Fprintln(os.Stderr, "xato:", name, "katalogda yoʻq")
		os.Exit(1)
	}
	fmt.Println(name, price/tiyinPerSom, "soʻm")
}
