package main

import "fmt"

func main() {
	stock := 40
	defer fmt.Println("defer: qoldiq", stock)
	stock -= 3
	fmt.Println("hozir: qoldiq", stock)
}
