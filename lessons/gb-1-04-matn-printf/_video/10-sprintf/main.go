package main

import "fmt"

func main() {
	name := "choy"
	sold := 7
	item := fmt.Sprintf("%s, %d dona", name, sold)
	fmt.Println(item)
}
