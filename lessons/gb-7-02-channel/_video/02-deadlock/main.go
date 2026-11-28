package main

import "fmt"

func main() {
	counts := make(chan int)
	counts <- 40
	n := <-counts
	fmt.Println(n)
}
