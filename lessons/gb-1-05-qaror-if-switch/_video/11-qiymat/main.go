package main

import "fmt"

func main() {
	name := "oʻrik"
	switch name {
	case "olma", "oʻrik":
		fmt.Println("meva")
	case "non":
		fmt.Println("non mahsuloti")
	default:
		fmt.Println("boshqa")
	}
}
