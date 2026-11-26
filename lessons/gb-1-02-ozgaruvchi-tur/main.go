package main

import "fmt"

func main() {
	fmt.Println("bozor ishga tushdi")
	fmt.Println("versiya: 0.1.0")

	var name string
	var price int
	var active bool
	fmt.Println(name, price, active)

	name = "olma"
	price = 12000
	active = true
	stock := 40
	weight := 1.5
	fmt.Println(name, price, active, stock, weight)
}
