package main

import (
	"fmt"
	"unicode/utf8"
)

// Javob: "Gʻayrat" = 8 bayt, 7 rune (G ʻ a y r a t; ʻ ikki bayt).
// "G'ayrat" (klaviatura apostrofi) = 7 bayt, 7 rune: faqat bayt soni oʻzgaradi.
func main() {
	for _, ism := range []string{"Gʻayrat", "G'ayrat"} {
		fmt.Println(ism, len(ism), utf8.RuneCountInString(ism))
	}
}
