package main

import "fmt"

type Holat int

const ( // bazaga saqlanadigan raqamlar qoʻlda yoziladi va hech qachon oʻzgarmaydi
	Nomalum  Holat = 0
	Yangi    Holat = 1
	Tolangan Holat = 2
	Bekor    Holat = 3 // yangi holat roʻyxat oxiriga
)

func main() {
	var saqlangan Holat = 2
	fmt.Println(saqlangan == Tolangan, saqlangan == Bekor)
}
