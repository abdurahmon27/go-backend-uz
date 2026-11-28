package main

import "fmt"

func main() {
	counts := make(chan int, 2)
	counts <- 40
	close(counts)
	fmt.Println(<-counts)
	n, ok := <-counts
	fmt.Println(n, ok)
}
