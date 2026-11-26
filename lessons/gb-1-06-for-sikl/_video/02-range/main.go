package main

import "fmt"

func main() {
	qty := 3
	for i := range qty {
		fmt.Println(i)
	}
}
