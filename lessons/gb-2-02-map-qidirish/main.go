package main

import (
	"fmt"
	"os"
	"slices"
)

const tiyinPerSom = 100 // 1 soʻm = 100 tiyin

// reserve buyurtmadagi mahsulotlarni ombordan ayiradi. Vazifa: TASK.md.
func reserve(stock map[string]int, order []string) error {
	return nil
}

func main() {
	prices := map[string]int64{
		"olma": 12_000 * tiyinPerSom,
		"non":  4_000 * tiyinPerSom,
		"choy": 25_000 * tiyinPerSom,
	}
	order := os.Args[1:]
	count := make(map[string]int)
	var total int64
	for _, name := range order {
		price, ok := prices[name]
		if !ok {
			fmt.Fprintln(os.Stderr, "xato:", name, "katalogda yoʻq")
			os.Exit(1)
		}
		count[name]++
		total += price
	}
	var names []string
	for name := range count {
		names = append(names, name)
	}
	slices.Sort(names)
	for _, name := range names {
		fmt.Printf("%-5s %d dona\n", name, count[name])
	}
	fmt.Printf("jami: %d soʻm\n", total/tiyinPerSom)
}
