package main

import (
	"fmt"
	"unicode/utf8"
)

// Bitta soʻz, toʻrt xil apostrof. Koʻzga deyarli bir xil, baytlari har xil.
func main() {
	sozlar := []string{
		"o'zbek", // U+0027 APOSTROPHE (klaviaturadagi)
		"oʻzbek", // U+02BB MODIFIER LETTER TURNED COMMA (teskari vergul, toʻgʻrisi)
		"o’zbek", // U+2019 RIGHT SINGLE QUOTATION MARK (oʻng jingalak)
		"o‘zbek", // U+2018 LEFT SINGLE QUOTATION MARK (chap jingalak)
	}
	for _, s := range sozlar {
		fmt.Printf("%s  len=%d  rune=%d\n", s, len(s), utf8.RuneCountInString(s))
	}
}
