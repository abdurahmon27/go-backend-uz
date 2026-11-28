package main

import "fmt"

func main() {
	counts := make(chan int, 2)
	counts <- 40
	counts <- 120
	fmt.Println(len(counts), cap(counts))
	counts <- 5
}
