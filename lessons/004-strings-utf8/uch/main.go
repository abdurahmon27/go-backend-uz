package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	sozlar := []string{
		"o'zbek", // U+0027 APOSTROPHE (klaviaturadagi)
		"oʻzbek", // U+02BB MODIFIER LETTER TURNED COMMA
		"o’zbek", // U+2019 RIGHT SINGLE QUOTATION MARK
		"o‘zbek", // U+2018 LEFT SINGLE QUOTATION MARK
	}
	for _, s := range sozlar {
		fmt.Printf("%s  len=%d  rune=%d  % x\n", s, len(s), utf8.RuneCountInString(s), s)
	}
}
