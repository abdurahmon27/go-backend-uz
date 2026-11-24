package main

import "fmt"

func main() {
	for i := 0; i < 3; i++ { // 1. uch qismli
		fmt.Print(i, " ")
	}
	n := 1
	for n < 100 { // 2. faqat shart (while)
		n *= 2
	}
	fmt.Println("|", n)
	for { // 3. cheksiz
		n -= 50
		if n < 0 {
			break
		}
	}
	fmt.Println(n)
	for i := range 3 { // 4. son ustida (Go 1.22+)
		fmt.Print(i, " ")
	}
	fmt.Println()
}
