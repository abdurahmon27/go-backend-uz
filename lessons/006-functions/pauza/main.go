package main

import "fmt"

func main() {
	x := 1
	f := func() int { return x * 10 }
	x = 5
	fmt.Println(f())
}
