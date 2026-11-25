// toza: ../hid dagi uchta funksiya type parameter'siz. Natija bir xil.
// go run ./toza
package main

import (
	"fmt"
	"strconv"
)

type Buyurtma struct{ Summa int64 }

func (b Buyurtma) String() string { return fmt.Sprintf("buyurtma: %d soʻm", b.Summa) }

// 1: faqat metod kerak — interface.
func Chiqar(x fmt.Stringer) { fmt.Println(x.String()) }

// 2: any yetarli.
func Yoz(x any) { fmt.Println("log:", x) }

// 3: har bir tur uchun oʻz kodi — alohida funksiyalar.
func MatnSon(x int64) string   { return strconv.FormatInt(x, 10) }
func MatnSatr(x string) string { return x }

func main() {
	Chiqar(Buyurtma{150000})
	Yoz(42)
	fmt.Println(MatnSon(99000), MatnSatr("nok"))
}
