package main

import "fmt"

func main() {
	asos := []string{"non", "sut"}
	asos = append(asos, "suv")
	asos = asos[:len(asos):len(asos)] // yoki slices.Clip(asos)

	ali := append(asos, "tuxum")
	vali := append(asos, "choy")

	fmt.Println("ali: ", ali)
	fmt.Println("vali:", vali)
}
