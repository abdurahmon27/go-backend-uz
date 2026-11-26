// 🟡 Vazifa: bazaga yoki JSON'ga raqam emas, nom saqlang ("paid"), shunda roʻyxat tartibi muhim boʻlmaydi.
// 1) HolatNomi(h Holat) string — Holat → "unknown" / "new" / "paid" / "cancelled"
// 2) HolatNomdan(s string) (Holat, error) — teskarisi; nomaʼlum nom uchun error qaytarsin (GO-X.03).
// Tekshirish: go run uyga/orta/main.go
// Kutilgan natija:
//
//	paid 2 <nil>
//	0 nomaʼlum holat: "refund"
package main

import (
	"errors"
	"fmt"
)

type Holat int

const (
	Nomalum Holat = iota
	Yangi
	Tolangan
	Bekor
)

func HolatNomi(h Holat) string {
	// VAZIFA
	return ""
}

func HolatNomdan(s string) (Holat, error) {
	// VAZIFA
	return Nomalum, errors.New("hali yozilmagan")
}

func main() {
	h, err := HolatNomdan(HolatNomi(Tolangan))
	fmt.Println(HolatNomi(h), h, err)

	h, err = HolatNomdan("refund")
	fmt.Println(h, err)
}
