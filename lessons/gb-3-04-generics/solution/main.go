// bozor: bir xil kodni har xil tur uchun bitta generic funksiya bilan yozamiz (GB-3.04).
//
//	go run .
package main

import (
	"fmt"
	"io"
	"os"
)

// Product — katalogdagi bitta mahsulot (GB-2.03).
type Product struct {
	ID    int64
	Name  string
	Price int64 // tiyin
	Stock int
}

// OrderItem — buyurtmadagi bitta qator (GB-2.03).
type OrderItem struct {
	ProductID int64
	Price     int64 // tiyin, buyurtma paytidagi narx
	Quantity  int
}

// OrderStatus — buyurtma holati (GB-2.03).
type OrderStatus int

const (
	StatusNew OrderStatus = iota
	StatusPaid
)

// Order — bitta buyurtma (GB-2.03).
type Order struct {
	ID     int64
	Items  []OrderItem
	Status OrderStatus
}

// Total buyurtma summasini tiyinda qaytaradi (GB-2.05).
func (o *Order) Total() int64 {
	var total int64
	for _, it := range o.Items {
		total += it.Price * int64(it.Quantity)
	}
	return total
}

// Filter items ichidan keep true qaytargan elementlarni yangi slice'da qaytaradi;
// hech biri mos kelmasa nil (len 0). T — tur parametri: chaqiruvchi beradigan istalgan tur (any).
func Filter[T any](items []T, keep func(T) bool) []T {
	var out []T
	for _, x := range items {
		if keep(x) {
			out = append(out, x)
		}
	}
	return out
}

// Number — Sum qabul qiladigan turlar: faqat int va int64.
type Number interface {
	int | int64
}

// Sum nums'dagi sonlarning yigʻindisini qaytaradi; boʻsh slice uchun 0.
func Sum[N Number](nums []N) N {
	var total N
	for _, n := range nums {
		total += n
	}
	return total
}

// Contains items ichida x bormi. == faqat comparable turlarda bor,
// shuning uchun constraint any emas, comparable.
func Contains[T comparable](items []T, x T) bool {
	for _, v := range items {
		if v == x {
			return true
		}
	}
	return false
}

// writeTotal summani w'ga yozadi. w faqat Write metodi uchun kerak,
// shuning uchun parametr oddiy interface: tur parametri hech narsa bermaydi.
func writeTotal(w io.Writer, total int64) error {
	_, err := fmt.Fprintf(w, "jami: %d tiyin\n", total)
	return err
}

// SumBy har bir elementdan value bergan sonni qoʻshib chiqadi.
// Masalan SumBy(products, func(p Product) int { return p.Stock }) — jami qoldiq.
// N'ning constraint'i Number: + faqat shunda ruxsat (any bilan kompilyator rad etadi).
func SumBy[T any, N Number](items []T, value func(T) N) N {
	var total N
	for _, x := range items {
		total += value(x)
	}
	return total
}

// paidTotal toʻlangan buyurtmalar summasini tiyinda qaytaradi.
func paidTotal(orders []Order) int64 {
	paid := Filter(orders, func(o Order) bool { return o.Status == StatusPaid })
	return SumBy(paid, func(o Order) int64 { return o.Total() })
}

func main() {
	products := []Product{
		{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
		{ID: 2, Name: "non", Price: 400_000, Stock: 3},
		{ID: 3, Name: "choy", Price: 2_500_000, Stock: 0},
	}
	orders := []Order{
		{ID: 7, Status: StatusPaid, Items: []OrderItem{{ProductID: 1, Price: 1_200_000, Quantity: 3}}},
		{ID: 8, Status: StatusNew, Items: []OrderItem{{ProductID: 2, Price: 400_000, Quantity: 5}}},
		{ID: 9, Status: StatusPaid, Items: []OrderItem{{ProductID: 2, Price: 400_000, Quantity: 2}}},
	}

	low := Filter(products, func(p Product) bool { return p.Stock < 5 })
	for _, p := range low {
		fmt.Println("kam qolgan:", p.Name, p.Stock)
	}
	paid := Filter(orders, func(o Order) bool { return o.Status == StatusPaid })
	for _, o := range paid {
		fmt.Println("toʻlangan buyurtma:", o.ID)
	}

	fmt.Println(Sum([]int{40, 3, 0}), Sum([]int64{3_600_000, 800_000}))
	fmt.Println(Contains([]string{"olma", "non"}, "non"), Contains([]int64{1, 2}, 3))
	if err := writeTotal(os.Stdout, paidTotal(orders)); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}
