package main

import (
	"fmt"
	"strings"
)

func main() {
	saytda := "Oʻzbekiston"
	qidiruv := "O'zbekiston"
	// Faqat o/g dan keyingi apostrof (klaviatura ' va egri ’ ‘) → ʻ (U+02BB). Boshqa joydagi ' tegilmaydi:
	// maʼno dagi tutuq belgisi (ʼ U+02BC) teskari vergul emas. Toʻliqroq qoida va chegaralar: ../tozala/
	r := strings.NewReplacer(
		"o'", "oʻ", "g'", "gʻ", "O'", "Oʻ", "G'", "Gʻ",
		"o’", "oʻ", "g’", "gʻ", "O’", "Oʻ", "G’", "Gʻ",
		"o‘", "oʻ", "g‘", "gʻ", "O‘", "Oʻ", "G‘", "Gʻ",
	)
	fmt.Println(r.Replace(qidiruv) == saytda)
	fmt.Println(r.Replace("ma'no"))
	// egri apostroflar ham: true true
	fmt.Println(r.Replace("O’zbekiston") == saytda, r.Replace("O‘zbekiston") == saytda)
}
