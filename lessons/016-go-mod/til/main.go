package main

import "fmt"

func main() {
	for i := range 3 { // Go 1.22 dan beri: butun son ustida range
		fmt.Print(i, " ")
	}
	fmt.Println()
}
