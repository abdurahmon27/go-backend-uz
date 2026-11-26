package main

import "fmt"

func main() {
	const yuz = 100       // untyped constant: turi hali yoʻq
	var a int8 = yuz      // bu yerda int8 boʻldi
	var b float64 = yuz   // bu yerda float64 boʻldi
	m := yuz              // tur aytilmagan joy → default type: int (butun son konstantasi uchun)
	fmt.Println(a/8, b/8) // int8: kasr tashlanadi, float64: qoladi
	fmt.Printf("%T\n", m) // int
}
