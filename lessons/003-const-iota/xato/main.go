package main

import "fmt"

type Holat int

const (
	Nomalum Holat = iota
	Yangi
	Bekor // yangi holat oʻrtaga qoʻshildi
	Tolangan
)

func main() {
	var saqlangan Holat = 2 // kecha bazaga yozilgan: oʻshanda 2 = Tolangan edi
	fmt.Println(saqlangan == Tolangan, saqlangan == Bekor)
}
