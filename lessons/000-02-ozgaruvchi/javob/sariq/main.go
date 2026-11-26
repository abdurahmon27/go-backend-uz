package main

import "fmt"

func main() {
	a := 3
	b := 8
	vaqtincha := a
	a = b
	b = vaqtincha
	fmt.Println(a, b)
}
