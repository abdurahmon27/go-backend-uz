package main

import "fmt"

func main() {
	id := 65
	key := "user:" + string(id)
	fmt.Println(key)
}
