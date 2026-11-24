package main

import (
	"fmt"
	"os"
)

func main() {
	ism := os.Getenv("ISM")
	if ism == "" {
		ism = "dunyo"
	}
	fmt.Println("Salom,", ism)
}
