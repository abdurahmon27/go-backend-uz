package main

import "fmt"

func main() {
	// Konstanta ifoda aniq hisoblanadi: oraliqdagi 150 hech qayerga saqlanmaydi,
	// int8 ga faqat yakuniy 90 qo'yiladi.
	var a int8 = 100 + 50 - 60
	fmt.Println(a)
}
