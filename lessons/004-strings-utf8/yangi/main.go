package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	nfc := "\u00f6zbek"  // ö bitta belgi (U+00F6)
	nfd := "o\u0308zbek" // o + ustidagi ikki nuqta (U+0308)
	for _, s := range []string{nfc, nfd} {
		fmt.Printf("%s  len=%d  rune=%d\n", s, len(s), utf8.RuneCountInString(s))
	}
	fmt.Println(nfc == nfd)
}
