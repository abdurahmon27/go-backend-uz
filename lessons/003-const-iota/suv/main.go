package main

import "fmt"

func main() {
	const yuz = 100       // turi yoʻq (untyped) konstanta
	var a int8 = yuz      // int8 katakka tushdi
	var b float64 = yuz   // float64 katakka tushdi
	fmt.Println(a/8, b/8) // int8: kasr tashlanadi, float64: qoladi
}
