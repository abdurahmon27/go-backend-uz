// Oson: ismingizda nechta bayt va nechta rune bor?
package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	for _, ism := range []string{"Aziz", "Gʻayrat", "Oʻgʻiloy"} {
		fmt.Printf("%-10s bayt=%d rune=%d\n", ism, len(ism), utf8.RuneCountInString(ism))
	}
}
