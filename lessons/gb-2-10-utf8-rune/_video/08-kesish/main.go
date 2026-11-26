package main

import (
	"fmt"
	"unicode/utf8"
)

// short nomni chek ustuniga sigʻdiradi: limit'dan uzun boʻlsa kesadi.
func short(name string, limit int) string {
	if len(name) <= limit {
		return name
	}
	return name[:limit]
}

func main() {
	cut := short("oʻrik", 2)
	fmt.Printf("%q\n", cut)
	fmt.Println(utf8.ValidString(cut))
}
