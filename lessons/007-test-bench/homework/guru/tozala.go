package tozala

import (
	"strings"
	"unicode"
)

// tozala: 004-darsdagi qoida, lekin hech narsa oʻzgarmasa,
// yangi string yaratmasdan s ning oʻzini qaytaradi.
func tozala(s string) string {
	r := []rune(s)
	ozgardi := false
	for i := 1; i < len(r); i++ {
		if !strings.ContainsRune("'`´‘’ʻʼ", r[i]) || !unicode.IsLetter(r[i-1]) {
			continue
		}
		ogdan := strings.ContainsRune("oOgG", r[i-1])
		keyingiHarf := i+1 < len(r) && unicode.IsLetter(r[i+1])
		soz2 := i >= 2 && unicode.IsLetter(r[i-2])
		yangi := r[i]
		switch {
		case ogdan && (keyingiHarf || soz2):
			yangi = 'ʻ'
		case keyingiHarf:
			yangi = 'ʼ'
		}
		if yangi != r[i] {
			r[i] = yangi
			ozgardi = true
		}
	}
	if !ozgardi {
		return s
	}
	return string(r)
}
