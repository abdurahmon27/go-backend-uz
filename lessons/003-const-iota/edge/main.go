package main

import "fmt"

const huge = 1 << 512

func main() {
	fmt.Println(huge >> 511)
}
