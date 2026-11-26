package main

import "fmt"

func main() {
	const ( // faylga saqlanadigan raqamlar qoʻlda yoziladi va hech qachon oʻzgarmaydi
		Nomalum  = 0
		Yangi    = 1
		Tolangan = 2
		Bekor    = 3 // yangi holat roʻyxat oxiriga
	)
	saqlangan := 2
	fmt.Println(saqlangan == Tolangan, saqlangan == Bekor)
}
