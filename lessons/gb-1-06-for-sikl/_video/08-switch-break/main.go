package main

import "fmt"

func main() {
	qty := 4
	for i := 1; i <= qty; i++ {
		switch {
		case i%3 == 0:
			break
		}
		fmt.Println(i)
	}
}
