package main

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

// short nomni limit ta rune'gacha qisqartiradi va rune'ni hech qachon oʻrtasidan kesmaydi.
func short(name string, limit int) string {
	n := 0
	for i := range name {
		if n == limit {
			return name[:i]
		}
		n++
	}
	return name
}

// shortLabel nom width ta rune'ga sigʻsa, oʻzini qaytaradi; sigʻmasa, birinchi width-1 ta rune
// va oxirida «…» (natija jami width ta rune; «…» bitta rune, lekin uch bayt).
func shortLabel(name string, width int) string {
	if utf8.RuneCountInString(name) <= width {
		return name
	}
	return short(name, width-1) + "…"
}

// maxNameRunes — mahsulot nomining eng katta uzunligi, rune'larda.
const maxNameRunes = 20

// checkName boʻsh nom, UTF-8 boʻlmagan nom va maxNameRunes'dan uzun nom uchun xato qaytaradi.
func checkName(name string) error {
	if name == "" {
		return errors.New("nom boʻsh")
	}
	if !utf8.ValidString(name) {
		return errors.New("nom UTF-8 emas")
	}
	if n := utf8.RuneCountInString(name); n > maxNameRunes {
		return fmt.Errorf("nom %d belgi, eng koʻpi %d", n, maxNameRunes)
	}
	return nil
}

func main() {
	names := []string{"olma", "oʻrik qoqi", "paxta yogʻi"}
	for _, name := range names {
		fmt.Printf("[%-10s]\n", short(name, 10))
	}
}
