package main

import "fmt"

func main() {
	var stock uint = 3 // omborda 3 ta telefon
	order := uint(5)   // mijoz 5 ta buyurtma qildi

	stock -= order
	fmt.Println("qoldiq:", stock)
}
