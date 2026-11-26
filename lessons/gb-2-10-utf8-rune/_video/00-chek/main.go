package main

import "fmt"

// short nomni chek ustuniga sigʻdiradi: limit'dan uzun boʻlsa kesadi.
func short(name string, limit int) string {
	if len(name) <= limit {
		return name
	}
	return name[:limit]
}

func main() {
	names := []string{"olma", "oʻrik qoqi", "paxta yogʻi"}
	for _, name := range names {
		fmt.Printf("[%-10s]\n", short(name, 10))
	}
}
