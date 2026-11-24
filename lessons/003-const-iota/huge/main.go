package main

import "fmt"

const huge = 1 << 511

func main() {
	fmt.Println(huge >> 510)
}
