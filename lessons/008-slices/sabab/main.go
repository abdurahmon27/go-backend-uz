package main

import (
	"fmt"
	"unsafe"
)

func main() {
	asos := []string{"non", "sut"}
	asos = append(asos, "suv")
	fmt.Println("len:", len(asos), "cap:", cap(asos))

	ali := append(asos, "tuxum")
	vali := append(asos, "choy")
	fmt.Println("bitta massivmi?", &ali[0] == &vali[0])
	fmt.Println("slice oʻzi:", unsafe.Sizeof(ali), "bayt")
}
