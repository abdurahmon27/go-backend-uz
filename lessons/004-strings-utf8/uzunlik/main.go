package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "oʻzbek"
	fmt.Println(len(s))                    // baytlar soni
	fmt.Println(utf8.RuneCountInString(s)) // rune'lar (belgilar) soni
}
