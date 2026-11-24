package main

import "fmt"

func main() {
	domenlar := []string{"a.uz", "b.uz", "c.uz"}
	var navbat []*string
	for _, d := range domenlar {
		navbat = append(navbat, &d) // d'ning manzili
	}
	for _, p := range navbat {
		fmt.Print(*p, " ")
	}
	fmt.Println()
}
