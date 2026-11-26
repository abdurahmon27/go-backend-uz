package main

import (
	"fmt"
	"os"
	"slices"
)

const tiyinPerSom = 100 // 1 soʻm = 100 tiyin

// reserve buyurtmadagi mahsulotlarni ombordan ayiradi. Biror mahsulot katalogda
// yoʻq boʻlsa yoki yetmasa, xato qaytaradi va omborni oʻzgartirmaydi.
func reserve(stock map[string]int, order []string) error {
	need := make(map[string]int)
	for _, name := range order {
		need[name]++
	}
	// avval hammasini tekshiramiz: shu yerda ombor hali oʻzgarmagan
	for _, name := range order {
		have, ok := stock[name]
		if !ok {
			return fmt.Errorf("%s omborda yoʻq", name)
		}
		if have < need[name] {
			return fmt.Errorf("%s: omborda %d dona, buyurtmada %d", name, have, need[name])
		}
	}
	// keyin ayiramiz: stock chaqiruvchining jadvali, oʻzgarish unga ham koʻrinadi
	for name, n := range need {
		stock[name] -= n
	}
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
