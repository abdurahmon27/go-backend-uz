// literal: struct e'lon qilish, composite literal, zero value va embedding.
package main

import "fmt"

type Manzil struct {
	Shahar string
	Kocha  string
}

type Mijoz struct {
	Ism    string
	Balans int64 // tiyinda
	Manzil       // embedding: tur nomi (Manzil) maydon nomi ham boʻladi
}

func main() {
	// kalitli literal: maydon nomi bilan, tartib muhim emas; Manzil yozilmadi -> zero value
	a := Mijoz{Ism: "Aziza", Balans: 50_000}
	fmt.Printf("%+v\n", a)

	// var: hamma maydon zero value oladi
	var b Mijoz
	fmt.Printf("%+v\n", b)

	// embedding: Manzil'ning maydonlari Mijoz'ga "koʻtariladi"
	a.Shahar = "Samarqand"
	fmt.Println(a.Shahar, a.Manzil.Shahar)
}
