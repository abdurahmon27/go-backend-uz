package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	price, err := strconv.Atoi(os.Args[3])
	fmt.Println(price, err)
}
