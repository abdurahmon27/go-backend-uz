package main

import "fmt"

// Oson: ishga tushirmasdan ayting, keyin tekshiring.
func main() {
	s := make([]int, 3, 5)
	t := append(s, 7)
	s[0] = 9
	fmt.Println(t[0], len(t), cap(t))
}
