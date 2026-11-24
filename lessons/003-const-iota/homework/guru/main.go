package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.January, int(time.January)) // January 1
	fmt.Println(time.Sunday, int(time.Sunday))   // Sunday 0
	var m time.Month
	fmt.Println(m) // %!Month(0): nol oy yo'q
}
