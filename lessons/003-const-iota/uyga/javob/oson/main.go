package main

import "fmt"

func main() {
	const (
		Nomalum  = iota // 0
		Yangi           // 1
		Tolangan        // 2: kechagi raqam joyida qoldi
		Bekor           // 3: oxiriga koʻchirildi
	)
	saqlangan := 2 // kecha faylga yozilgan: 2 = Tolangan
	fmt.Println(saqlangan == Tolangan, saqlangan == Bekor)
}
