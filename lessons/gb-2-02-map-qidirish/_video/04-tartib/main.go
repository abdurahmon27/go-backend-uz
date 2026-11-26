package main

import (
	"fmt"
	"os"
)

func main() {
	order := os.Args[1:]
	count := make(map[string]int)
	for _, name := range order {
		count[name]++
	}
	for name, n := range count {
		fmt.Println(name, n)
	}
}
