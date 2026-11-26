package main

import "fmt"

const tiyinPerSom = 100    // 1 soʻm = 100 tiyin
const vatPercent = 12      // QQS, foizda
const deliverySom = 15_000 // yetkazib berish, soʻmda (QQS qoʻshilmaydi)

func main() {
	fmt.Println("bozor ishga tushdi")
	fmt.Println("versiya: 0.1.0")

	var name string
	var price int64
	var active bool
	fmt.Println(name, price, active)

	name = "choy"
	price = 1_999 * tiyinPerSom // tiyinda
	active = true
	stock := 40
	weight := 1.5
	fmt.Println(name, price, active, stock, weight)

	sold := 7
	stock = stock - sold
	fmt.Println(name, stock)

	subtotal := price * int64(sold)
	vat := subtotal * vatPercent / 100
	total := subtotal + vat + deliverySom*tiyinPerSom
	fmt.Println(name, sold, "dona")
	fmt.Println("QQS:", vat/tiyinPerSom, "soʻm", vat%tiyinPerSom, "tiyin")
	fmt.Println("jami:", total/tiyinPerSom, "soʻm", total%tiyinPerSom, "tiyin")
}
