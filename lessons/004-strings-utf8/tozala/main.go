package main

import (
	"fmt"
	"strings"
	"unicode"
)

// tozala apostrof turlarini bir xil qiladi:
// o/g dan keyin -> ʻ (U+02BB), boshqa harfdan keyin -> ʼ (U+02BC, tutuq belgisi).
// Faqat harfdan keyin turganlarini oʻzgartiradi: soʻz boshidagi ‘ qoʻshtirnoqqa tegmaydi.
func tozala(s string) string {
	r := []rune(s)
	for i := 1; i < len(r); i++ {
		if !strings.ContainsRune("'`´‘’ʻʼ", r[i]) || !unicode.IsLetter(r[i-1]) {
			continue
		}
		ogdan := strings.ContainsRune("oOgG", r[i-1])
		keyingiHarf := i+1 < len(r) && unicode.IsLetter(r[i+1])
		soz2 := i >= 2 && unicode.IsLetter(r[i-2]) // o/g soʻz boshida emas
		switch {
		case ogdan && (keyingiHarf || soz2): // oʻzbek, togʻ
			r[i] = 'ʻ'
		case keyingiHarf: // maʼno
			r[i] = 'ʼ'
		}
	}
	return string(r)
}

func main() {
	fmt.Println(tozala("O'zbekiston, tog’, ma'no, O`zbek"))
	fmt.Println(tozala("U ‘Salom’ dedi"))
	baza := tozala("Bu togʻ baland")
	fmt.Println(strings.Contains(baza, tozala("tog'")))
}
