// oson: Indeks — x birinchi uchragan oʻrin yoki -1. Bitta generic funksiya.
// go run ./homework/oson
package main

import (
	"fmt"
	"slices"
)

func Indeks[T comparable](s []T, x T) int {
	for i, v := range s {
		if v == x {
			return i
		}
	}
	return -1
}

func main() {
	fmt.Println(Indeks([]int64{150000, 99000}, 99000))
	fmt.Println(Indeks([]string{"olma", "nok"}, "uzum"))
	// Standart kutubxonada tayyori bor (Go 1.21+):
	fmt.Println(slices.Index([]int64{150000, 99000}, 99000))
}
