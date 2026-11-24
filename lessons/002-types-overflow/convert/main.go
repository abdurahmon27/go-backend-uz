package main

import "fmt"

func main() {
	big := 300
	small := int8(big)
	fmt.Println(small)

	other := 200
	fmt.Println(int8(other))
}
