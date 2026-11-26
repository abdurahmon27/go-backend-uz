package main

import "fmt"

func main() {
	jami := 0
	for oy := 1; oy <= 12; oy++ {
		jami = jami + 300
		fmt.Println(oy, "oy:", jami)
		if jami >= 900 {
			break
		}
	}
}
