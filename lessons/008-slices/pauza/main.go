package main

import "fmt"

func main() {
	asos := []string{"non", "sut", "suv"}

	ali := append(asos, "tuxum")
	vali := append(asos, "choy")

	fmt.Println("ali: ", ali)
	fmt.Println("vali:", vali)
}
