package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	price, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
	fmt.Println("narx:", price)
}
