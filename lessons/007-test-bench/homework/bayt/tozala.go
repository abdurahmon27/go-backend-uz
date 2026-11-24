package tozala

import "strings"

// tozalaCA: tezfix'dagi tez yoʻl (strings.ContainsAny).
func tozalaCA(s string) string {
	if !strings.ContainsAny(s, "'`´‘’ʻʼ") {
		return s
	}
	return tozalaEski(s)
}

// tozala: tez yoʻl baytlar boʻyicha. ' va ` dan tashqari hamma apostroflar
// UTF-8'da 0x80 dan katta baytlardan boshlanadi: bunday bayt boʻlsa,
// sekin yoʻlga oʻtamiz (u oʻzi hamma narsani tekshiradi).
func tozala(s string) string {
	for i := 0; i < len(s); i++ {
		if c := s[i]; c == '\'' || c == '`' || c >= 0x80 {
			return tozalaEski(s)
		}
	}
	return s
}
