// Savol: sikl birinchi qatordayoq break qiladi. "fayl yopildi" chiqadimi? Qachon?
package main

import (
	"fmt"
	"iter"
)

func Qatorlar() iter.Seq[string] {
	return func(yield func(string) bool) {
		fmt.Println("fayl ochildi")
		defer fmt.Println("fayl yopildi")
		for _, q := range []string{"birinchi", "ikkinchi", "uchinchi"} {
			if !yield(q) {
				return
			}
		}
	}
}

func main() {
	for q := range Qatorlar() {
		fmt.Println(q)
		break
	}
	fmt.Println("sikldan keyin")
}
