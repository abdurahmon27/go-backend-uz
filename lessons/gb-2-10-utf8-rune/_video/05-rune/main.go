package main

import "fmt"

func main() {
	var r rune
	r = 'ʻ'
	fmt.Printf("%d %c %U\n", r, r, r)
}
