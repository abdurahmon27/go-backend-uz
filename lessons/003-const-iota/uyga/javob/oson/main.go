// 🟢 Javob: yangi holat faqat roʻyxat oxiriga qoʻshiladi, shunda eski raqamlar joyida qoladi.
// (Qoʻlda yozilgan raqamlar bilan ham boʻladi: Qaytarilgan Holat = 4.)
package main

import "fmt"

type Holat int

const (
	Nomalum     Holat = iota // 0
	Yangi                    // 1
	Tolangan                 // 2: kechagi raqam joyida
	Bekor                    // 3
	Qaytarilgan              // 4: yangi holat → oxiriga
)

func main() {
	var saqlangan Holat = 2
	fmt.Println(saqlangan == Tolangan, saqlangan == Bekor, saqlangan == Qaytarilgan)
}
