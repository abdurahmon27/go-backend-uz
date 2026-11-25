// Cheksiz rekursiya: stack bir milliard baytdan oshadi, fatal error.
package main

import "fmt"

func sana(n int) int {
	return sana(n+1) + 1
}

func main() {
	defer func() {
		r := recover()
		fmt.Println("qutqarildi:", r)
	}()

	fmt.Println(sana(0))
}
