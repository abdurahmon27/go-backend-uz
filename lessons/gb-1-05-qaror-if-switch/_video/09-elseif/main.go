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

	name = "oʻrik"
	price = 25_000 * tiyinPerSom // tiyinda
	active = true
	stock := 40
	weight := 1.5
	fmt.Println(name, price, active, stock, weight)

	sold := 24
	if sold <= 0 {
		fmt.Println("xato: miqdor kamida 1 dona boʻlishi kerak")
		return
	}
	if left := stock - sold; left < 0 {
		fmt.Printf("xato: omborda %d dona yetmaydi\n", -left)
		return
	}
	stock = stock - sold
	fmt.Println(name, stock)

	subtotal := price * int64(sold)
	vat := subtotal * vatPercent / 100

	var delivery int64 // tiyinda
	if subtotal >= 500_000*tiyinPerSom {
		delivery = 0
	} else if subtotal >= 100_000*tiyinPerSom {
		delivery = 10_000 * tiyinPerSom
	} else {
		delivery = deliverySom * tiyinPerSom
	}
	total := subtotal + vat + delivery

	item := fmt.Sprintf("%s, %d dona", name, sold)
	fmt.Printf("%-14s %9d.%02d\n", item, subtotal/tiyinPerSom, subtotal%tiyinPerSom)
	fmt.Printf("%-14s %9d.%02d\n", "QQS", vat/tiyinPerSom, vat%tiyinPerSom)
	fmt.Printf("%-14s %9d.%02d\n", "yetkazish", delivery/tiyinPerSom, delivery%tiyinPerSom)
	fmt.Printf("%-14s %9d.%02d\n", "jami", total/tiyinPerSom, total%tiyinPerSom)
}
