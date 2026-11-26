package main

import (
	"fmt"
	"os"
	"slices"
)

func main() {
	order := os.Args[1:]
	count := make(map[string]int)
	for _, name := range order {
		count[name]++
	}
	var names []string
	for name := range count {
		names = append(names, name)
	}
	slices.Sort(names)
	for _, name := range names {
		fmt.Printf("%-5s %d dona\n", name, count[name])
	}
}
