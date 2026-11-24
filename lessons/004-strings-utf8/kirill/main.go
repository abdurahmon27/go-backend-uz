package main

import "fmt"

func main() {
	lotin := "A"
	kirill := "\u0410" // lex.uz alifbo jadvalidagi birinchi "A"
	for _, s := range []string{lotin, kirill} {
		fmt.Printf("%s  %U  len=%d\n", s, []rune(s)[0], len(s))
	}
	fmt.Println(lotin == kirill)
}
