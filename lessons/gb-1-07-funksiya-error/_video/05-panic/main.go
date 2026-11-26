package main

import (
	"fmt"
	"os"
	"strconv"
)

const tiyinPerSom = 100    // 1 soʻm = 100 tiyin
const vatPercent = 12      // QQS, foizda
const deliverySom = 15_000 // yetkazib berish, soʻmda (QQS qoʻshilmaydi)
const maxFree = 2          // bitta buyurtmada koʻpi bilan shuncha dona bepul

// Buyurtma. Nom va narx endi buyruq qatoridan keladi (main ichida).
var active = true // sotuvda
var stock = 40    // omborda, dona
var sold = 7      // buyurtmada, dona; vazifada buyruq qatoridan

// parsePrice soʻmdagi narx matnini tiyinga oʻgiradi.
func parsePrice(s string) (int64, error) {
	som, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if som <= 0 {
		return 0, fmt.Errorf("narx musbat boʻlishi kerak, keldi %d", som)
	}
	return int64(som) * tiyinPerSom, nil
}

func main() {
	fmt.Println("bozor ishga tushdi")
	fmt.Println("versiya: 0.1.0")

	name := os.Args[2]
	price, err := parsePrice(os.Args[3])
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}

	if name == "" {
		fmt.Println("xato: mahsulot nomi boʻsh")
		return
	}
	if !active {
		fmt.Printf("xato: %s sotuvda emas\n", name)
		return
	}
	if price > 10_000_000*tiyinPerSom { // price <= 0 endi parsePrice ichida
		fmt.Printf("xato: narx notoʻgʻri: %d tiyin\n", price)
		return
	}
	if sold <= 0 {
		fmt.Println("xato: miqdor kamida 1 dona boʻlishi kerak")
		return
	}
	if sold > stock {
		fmt.Printf("xato: omborda %d dona bor, %d soʻraldi\n", stock, sold)
		return
	}
	if left := stock - sold; left < 5 {
		fmt.Printf("diqqat: omborda %d dona qoladi\n", left)
	}

	var subtotal int64
	free := 0
	for i := 1; i <= sold; i++ {
		if i%3 == 0 && free < maxFree {
			free++
			continue // aksiya: har uchinchi dona bepul
		}
		subtotal += price
	}
	vat := subtotal * vatPercent / 100

	var delivery int64 // tiyinda
	switch {
	case subtotal >= 500_000*tiyinPerSom:
		delivery = 0
	case subtotal >= 100_000*tiyinPerSom:
		delivery = 10_000 * tiyinPerSom
	default:
		delivery = deliverySom * tiyinPerSom
	}
	total := subtotal + vat + delivery

	item := fmt.Sprintf("%s, %d dona", name, sold)
	fmt.Printf("%-14s %9d.%02d\n", item, subtotal/tiyinPerSom, subtotal%tiyinPerSom)
	fmt.Printf("%-14s %9d.%02d\n", "QQS", vat/tiyinPerSom, vat%tiyinPerSom)
	fmt.Printf("%-14s %9d.%02d\n", "yetkazish", delivery/tiyinPerSom, delivery%tiyinPerSom)
	fmt.Printf("%-14s %9d.%02d\n", "jami", total/tiyinPerSom, total%tiyinPerSom)
	if free > 0 {
		fmt.Printf("aksiya: %d dona bepul\n", free)
	}
}
