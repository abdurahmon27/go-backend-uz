package main

import (
	"fmt"
	"os"
)

func main() {
	ism, bor := os.LookupEnv("ISM")
	switch {
	case !bor:
		fmt.Println("ISM berilmagan: Salom, dunyo")
	case ism == "":
		fmt.Println("ISM bor, lekin boʻsh: Salom, kim?")
	default:
		fmt.Println("Salom,", ism)
	}
}
