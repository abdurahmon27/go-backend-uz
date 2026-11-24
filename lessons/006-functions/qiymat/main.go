package main

import (
	"fmt"
	"strings"
)

func hammasiga(sozlar []string, f func(string) string) []string {
	var natija []string
	for _, s := range sozlar {
		natija = append(natija, f(s))
	}
	return natija
}

func main() {
	sozlar := []string{"olma", "nok", "uzum"}
	fmt.Println(hammasiga(sozlar, strings.ToUpper))

	qavs := func(s string) string { return "[" + s + "]" }
	fmt.Println(hammasiga(sozlar, qavs))
	fmt.Printf("%T\n", qavs)
}
