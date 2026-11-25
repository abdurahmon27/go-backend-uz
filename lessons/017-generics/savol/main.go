// savol: Katta(3, 2.5) — kompilyatsiya boʻladimi? T nima boʻladi?
// go run ./savol
package main

import (
	"cmp"
	"fmt"
)

func Katta[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func main() {
	k := Katta(3, 2.5)
	fmt.Printf("%v %T\n", k, k)
}
