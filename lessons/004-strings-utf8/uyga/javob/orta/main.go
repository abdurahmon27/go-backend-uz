package main

import (
	"fmt"
	"unicode/utf8"
)

// Javob: len baytlarni sanaydi: "Gʻani" = 6 bayt (ʻ ikki bayt), lekin 5 rune.
// Belgilarni sanash uchun utf8.RuneCountInString.
func main() {
	ism := "Gʻani"
	if utf8.RuneCountInString(ism) <= 5 {
		fmt.Println("qabul")
	} else {
		fmt.Println("juda uzun")
	}
}
