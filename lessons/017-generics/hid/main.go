// hid: uchta ortiqcha generic. Hammasi kompilyatsiya boʻladi va ishlaydi,
// lekin type parameter hech narsa bermaydi. Tozalangani: ../toza
// go run ./hid
package main

import (
	"fmt"
	"strconv"
)

type Buyurtma struct{ Summa int64 }

func (b Buyurtma) String() string { return fmt.Sprintf("buyurtma: %d soʻm", b.Summa) }

// 1-hid: T bilan faqat metod chaqiriladi.
func Chiqar[T fmt.Stringer](x T) { fmt.Println(x.String()) }

// 2-hid: T hech narsa bermaydi, Println baribir any oladi.
func Yoz[T any](x T) { fmt.Println("log:", x) }

// 3-hid: har bir tur uchun boshqa kod (type switch).
func Matn[T int64 | string](x T) string {
	switch v := any(x).(type) {
	case int64:
		return strconv.FormatInt(v, 10)
	case string:
		return v
	}
	return ""
}

func main() {
	Chiqar(Buyurtma{150000})
	Yoz(42)
	fmt.Println(Matn(int64(99000)), Matn("nok"))
}
