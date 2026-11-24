package main

import (
	"fmt"
	"strings"
)

func birlashtir(ajratgich string, sozlar ...string) string {
	var b strings.Builder
	for i, s := range sozlar {
		if i > 0 {
			b.WriteString(ajratgich)
		}
		b.WriteString(s)
	}
	return b.String()
}

func main() {
	fmt.Println(birlashtir(", ", "olma", "nok", "uzum"))
	fmt.Printf("%q\n", birlashtir("-"))
	meva := []string{"anor", "behi"}
	fmt.Println(birlashtir(" + ", meva...))
	fmt.Println(strings.Join(meva, " + ")) // standart kutubxonada tayyori bor
}
