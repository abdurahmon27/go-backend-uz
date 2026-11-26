package main

import "fmt"

type Holat int

const (
	Nomalum Holat = iota
	// Yangi oʻchirildi → keyingilarning raqami bittaga kamaydi
	Tolangan // endi 1
	Bekor    // endi 2
)

const (
	Nomalum2  Holat = iota
	_               // Yangi oʻchirildi, lekin raqami band qoldi
	Tolangan2       // 2
	Bekor2          // 3
)

func main() {
	var saqlangan Holat = 2 // kecha: 2 = Tolangan
	fmt.Println(saqlangan == Tolangan, saqlangan == Bekor)
	fmt.Println(saqlangan == Tolangan2, saqlangan == Bekor2)
}
