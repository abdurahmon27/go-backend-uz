package main

import "fmt"

func main() {
	// Bu dastur 1 + 2 + 3 + 4 + 5 = 15 ni chiqarishi kerak edi, lekin 5 chiqaradi.
	// Yashirin xatoni toping va tuzating.
	jami := 0
	for i := 1; i <= 5; i++ {
		jami = 0
		jami = jami + i
	}
	fmt.Println(jami)
}
