package main

import "fmt"

func main() {
	var ishlar []func()
	for i := 0; i < 3; i++ {
		ishlar = append(ishlar, func() { fmt.Print(i, " ") })
	}
	for _, ish := range ishlar {
		ish()
	}
	fmt.Println()
}
