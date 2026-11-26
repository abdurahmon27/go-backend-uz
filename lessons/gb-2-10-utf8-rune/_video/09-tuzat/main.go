package main

import "fmt"

// short nomni limit ta rune'gacha qisqartiradi va rune'ni hech qachon oʻrtasidan kesmaydi.
func short(name string, limit int) string {
	n := 0
	for i := range name {
		if n == limit {
			return name[:i]
		}
		n++
	}
	return name
}

func main() {
	names := []string{"olma", "oʻrik qoqi", "paxta yogʻi"}
	for _, name := range names {
		fmt.Printf("[%-10s]\n", short(name, 10))
	}
}
