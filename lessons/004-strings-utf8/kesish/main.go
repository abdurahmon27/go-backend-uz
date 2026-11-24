package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "oʻzbek"
	fmt.Printf("%q\n", s[:2])                 // 2 bayt: belgining yarmi
	fmt.Println(utf8.ValidString(s[:2]))      // bu UTF-8 emas
	fmt.Printf("%q\n", string([]rune(s)[:2])) // 2 rune: butun "oʻ"
	fmt.Printf("%q\n", string([]rune(s)[:1])) // 1 rune: faqat "o"
}
