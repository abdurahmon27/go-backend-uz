package main

import "fmt"

func main() {
	var son int16
	for i := 1; i <= 100000; i++ {
		son++
		if son < 0 {
			fmt.Println(i, son)
			break
		}
	}
}
