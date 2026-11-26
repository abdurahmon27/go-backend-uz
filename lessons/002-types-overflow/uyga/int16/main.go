package main

import "fmt"

func main() {
	var son int16
	for i := 1; i <= 100000; i++ {
		son++
		// Vazifa: son manfiy boʻlsa (son < 0), i va son ni chiqaring
		// va break bilan sikldan chiqing.
	}
	fmt.Println("sikl tugadi")
}
