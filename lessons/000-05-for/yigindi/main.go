package main

import "fmt"

func main() {
	jami := 0
	for i := 1; i <= 5; i++ {
		jami = jami + i
	}
	fmt.Println(jami)
}
