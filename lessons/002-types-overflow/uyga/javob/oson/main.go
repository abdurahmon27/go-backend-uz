package main

import "fmt"

// Javob: i <= 127 boʻlsa 127 chiqadi (hali sigʻadi),
// i <= 128 boʻlsa -128 chiqadi (128-laykda son doira boshiga aylanadi).
func main() {
	var layk int8
	for i := 1; i <= 128; i++ {
		layk++
	}
	fmt.Println(layk)
}
