package tozala

import (
	"strings"
	"unicode"
)

// tozalaEski: 004-darsdagi tozala, oʻzgarishsiz.
// U apostrof turlarini bir xil qiladi:
// o/g dan keyin -> ʻ (U+02BB), boshqa harfdan keyin -> ʼ (U+02BC, tutuq belgisi).
// Faqat harfdan keyin turganlarini oʻzgartiradi: soʻz boshidagi ‘ qoʻshtirnoqqa tegmaydi.
func tozalaEski(s string) string {
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
