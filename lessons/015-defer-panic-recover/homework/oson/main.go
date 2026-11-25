// Oson: avval natijani qog'ozda yozing, keyin ishga tushiring.
package main

import "fmt"

func main() {
	x := 1
	defer fmt.Println("A:", x)
	for i := range 3 {
		defer fmt.Println("i =", i)
	}
	defer func() { fmt.Println("B:", x) }()
	x = 2
}
