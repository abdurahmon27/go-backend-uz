package main

import "fmt"

func main() {
	s := "oʻzbek"
	for i, r := range s {
		fmt.Printf("%d  %c  %U  %d\n", i, r, r, r)
	}
}
