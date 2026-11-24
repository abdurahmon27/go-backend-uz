package main

import "fmt"

func main() {
	s := "oʻzbek"
	for i := 0; i < len(s); i++ {
		fmt.Printf("%c", s[i])
	}
	fmt.Println()
}
