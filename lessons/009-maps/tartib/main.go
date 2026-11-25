package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

func main() {
	narx := map[string]int{"non": 5000, "sut": 12000, "tuxum": 18000, "choy": 25000, "olma": 15000}

	const n = 100_000
	sanoq := map[string]int{} // har bir tartib necha marta chiqdi
	for range n {
		var kalitlar []string
		for nom := range narx {
			kalitlar = append(kalitlar, nom)
		}
		sanoq[strings.Join(kalitlar, " ")]++
	}
	fmt.Printf("%d ta sikl, %d xil tartib:\n", n, len(sanoq))
	for _, t := range slices.Sorted(maps.Keys(sanoq)) {
		fmt.Printf("%-24s %5.1f%%\n", t, 100*float64(sanoq[t])/n)
	}
}
