package main

import "fmt"

func main() {
	s := make([]int, 0, 4)
	s = append(s, 1)
	fmt.Println(len(s), cap(s))
}
