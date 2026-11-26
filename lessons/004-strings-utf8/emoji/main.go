package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "😀"                               // U+1F600 = 128512 > 65535
	fmt.Println(len(s))                    // baytlar
	fmt.Println(utf8.RuneCountInString(s)) // rune'lar
}
