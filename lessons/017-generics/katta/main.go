// katta: bir xil kod ikki marta, faqat tur boshqa → bitta generic funksiya.
// go run ./katta
package main

import (
	"cmp"
	"fmt"
)

func KattaInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func KattaString(a, b string) string {
	if a > b {
		return a
	}
	return b
}

// T — type parameter, cmp.Ordered — uning constraint'i.
func Katta[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func main() {
	fmt.Println(Katta(3, 7), Katta(1.5, 0.5), Katta("olma", "nok"))
	fmt.Println(Katta[int64](150000, 99000))
	fmt.Println(KattaInt(3, 7), KattaString("olma", "nok"))
}
