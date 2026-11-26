package main

import (
	"fmt"
	"os"
)

func main() {
	price := os.Args[3]
	fmt.Println(price + 1)
}
