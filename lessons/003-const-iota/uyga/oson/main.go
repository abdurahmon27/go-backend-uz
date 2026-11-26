// 🟢 Vazifa: pul qaytarilgan buyurtmalar uchun yangi holat — Qaytarilgan — qoʻshing.
// Shart: bazada allaqachon saqlangan raqamlar (2 = Tolangan, 3 = Bekor) maʼnosini saqlashi kerak.
// Avval natijani taxmin qiling, keyin: go run uyga/oson/main.go
// Kutilgan natija (toʻgʻri yechimda): true false false
package main

import "fmt"

type Holat int

const (
	Nomalum Holat = iota
	Yangi
	Tolangan
	Bekor
	// VAZIFA: Qaytarilgan holatini shunday qoʻshingki, yuqoridagi raqamlar oʻzgarmasin
)

func main() {
	var saqlangan Holat = 2 // kecha bazaga yozilgan: 2 = Tolangan
	fmt.Println(saqlangan == Tolangan, saqlangan == Bekor, saqlangan == Qaytarilgan)
}
