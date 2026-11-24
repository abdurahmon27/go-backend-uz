package main

import (
	"fmt"
	"strconv"
)

func soz(r rune) string {
	switch r {
	case '-':
		return "minus"
	case '0':
		return "nol"
	case '1':
		return "bir"
	case '2':
		return "ikki"
	case '3':
		return "uch"
	case '4':
		return "toʻrt"
	case '5':
		return "besh"
	case '6':
		return "olti"
	case '7':
		return "yetti"
	case '8':
		return "sakkiz"
	case '9':
		return "toʻqqiz"
	default:
		return "?"
	}
}

func main() {
	for _, n := range []int{2026, 1990, -7} {
		fmt.Print(n, ": ")
		for _, r := range strconv.Itoa(n) {
			fmt.Print(soz(r), " ")
		}
		fmt.Println()
	}
}
