package main

import "fmt"

func main() {
	counts := make(chan int)
	go func() {
		counts <- 40
	}()
	n := <-counts
	fmt.Println(n)
}
