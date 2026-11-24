package main

import (
	"fmt"
	"slices"
)

// mediana (soddalashtirilgan: toq sonli, boʻsh emas)
func mediana(vaqtlar ...int) int { // vaqtlar: []int
	slices.Sort(vaqtlar) // parametrni joyida saralaydi
	return vaqtlar[len(vaqtlar)/2]
}

func main() {
	fmt.Println(mediana(90, 10, 50))
	kechikish := []int{90, 10, 50}
	fmt.Println(mediana(kechikish...))
	fmt.Println(kechikish)
}
