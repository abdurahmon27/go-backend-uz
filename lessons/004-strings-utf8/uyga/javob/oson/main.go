package main

import (
	"fmt"
	"unicode/utf8"
)

// Javob: len baytlarni sanaydi: "Gʻani" = 6 bayt (ʻ ikki bayt), lekin 5 rune.
// Frontend'dagi .length UTF-16 birliklarini sanaydi; bu ismlarda u rune soniga teng, shuning uchun backend
// rune'larni sanashi kerak: utf8.RuneCountInString.
// (Emoji kabi U+FFFF dan katta belgilarda JS .length 2 beradi, rune esa 1 — bu yerda ular uchramaydi.)
func ismTekshir(ism string) string {
	if utf8.RuneCountInString(ism) <= 5 {
		return "qabul"
	}
	return "juda uzun"
}

func main() {
	for _, ism := range []string{"Anvar", "Gʻani", "Oʻktam"} {
		fmt.Println(ism, ismTekshir(ism))
	}
}
