package main

import "fmt"

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

// Vazifa (TASK.md, 1-qadam): nom width ta rune'ga sigʻsa, oʻzini qaytarsin; sigʻmasa,
// birinchi width-1 ta rune va oxirida «…» (natija jami width ta rune).
func shortLabel(name string, width int) string {
	return name
}

// Vazifa (TASK.md, 2-qadam): boʻsh nom, UTF-8 boʻlmagan nom va 20 rune'dan uzun nom uchun xato qaytarsin.
func checkName(name string) error {
	return nil
}

func main() {
	names := []string{"olma", "oʻrik qoqi", "paxta yogʻi"}
	for _, name := range names {
		fmt.Printf("[%-10s]\n", short(name, 10))
	}
}
