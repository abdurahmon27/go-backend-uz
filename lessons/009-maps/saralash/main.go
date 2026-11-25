package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	narx := map[string]int{"non": 5000, "sut": 12000, "tuxum": 18000, "choy": 25000, "olma": 15000}
	for _, nom := range slices.Sorted(maps.Keys(narx)) {
		fmt.Println(nom, narx[nom])
	}
}
