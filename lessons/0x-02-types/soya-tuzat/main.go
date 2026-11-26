package main

import "fmt"

func main() {
	total := 150
	discount := 0
	if total > 100 {
		discount = 10
		fmt.Println("chegirma:", discount)
	}
	fmt.Println(total - discount)
}
