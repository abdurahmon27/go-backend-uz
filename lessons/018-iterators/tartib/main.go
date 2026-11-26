// Tartib: map tartibsiz, maps.Keys esa slice emas — iterator.
package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	jami := map[string]int64{
		"Aziz":   150000,
		"Bobur":  99000,
		"Dilnoz": 42000,
		"Kamola": 310000,
	}

	for range 3 {
		var nomlar []string
		for k := range jami {
			nomlar = append(nomlar, k)
		}
		fmt.Println(nomlar)
	}

	fmt.Printf("%T\n", maps.Keys(jami))
	fmt.Println(slices.Sorted(maps.Keys(jami)))
}
