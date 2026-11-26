package main

import "fmt"

func main() {
	const a, b = 0.1, 0.2 // untyped: aniq hisoblanadi
	fmt.Println(a+b == 0.3)

	const c, d float64 = 0.1, 0.2 // typed: float64 ga yaxlitlangan
	fmt.Println(c+d == 0.3)

	x, y := 0.1, 0.2 // float64 oʻzgaruvchilar
	fmt.Println(x+y == 0.3, x+y)
}
