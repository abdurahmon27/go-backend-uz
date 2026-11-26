package main

import "fmt"

func main() {
	const (
		Nomalum = iota
		Yangi
		Bekor // VAZIFA: shu qatorni roʻyxat oxiriga koʻchiring
		Tolangan
	)
	saqlangan := 2 // kecha faylga yozilgan: oʻshanda 2 = Tolangan edi
	fmt.Println(saqlangan == Tolangan, saqlangan == Bekor)
}
