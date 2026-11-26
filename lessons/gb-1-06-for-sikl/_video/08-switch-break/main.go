package main

import "fmt"

func main() {
	sold := 4
	for i := 1; i <= sold; i++ {
		switch {
		case i%3 == 0:
			break
		}
		fmt.Println(i)
	}
}
