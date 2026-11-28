// bozor importeri: ikki yetkazib beruvchining narx roʻyxatini channel'lar orqali bosqichma-bosqich oʻtkazadi (GB-7.02).
//
//	parse (chorsu) ─┐
//	                ├─ raw ─ inStock ─ stocked ─ unique ─ ready ─ importAll
//	parse (oloy)  ──┘
//
// Javob (TASK.md): unique bosqichi yozilgan va importAll'ga qoʻyilgan.
package main

import (
	"fmt"
	"sync"
)

// parse qatorlarni Product'ga aylantirib out'ga yuboradi. out'ni yopmaydi: unga bir nechta parse yozadi,
// yopishni hammasini kutgan bitta goroutine qiladi (importAll).
func parse(lines []string, out chan<- Product) {
	for _, line := range lines {
		p, err := parseLine(line)
		if err != nil {
			continue // xatoni qaytarish: GB-7.04
		}
		out <- p
	}
}

// inStock faqat qoldigʻi bor mahsulotlarni oʻtkazadi; in yopilgach, out'ni yopadi.
func inStock(in <-chan Product, out chan<- Product) {
	defer close(out)
	for p := range in {
		if p.Stock > 0 {
			out <- p
		}
	}
}

// unique har ID'ni faqat birinchi marta oʻtkazadi; in yopilgach, out'ni yopadi.
func unique(in <-chan Product, out chan<- Product) {
	defer close(out)
	seen := make(map[int64]bool)
	for p := range in {
		if seen[p.ID] {
			continue
		}
		seen[p.ID] = true
		out <- p
	}
}

// importAll ikki roʻyxatni bir vaqtda oʻqiydi va qoldigʻi bor mahsulotlarni qaytaradi.
func importAll(a, b []string) []Product {
	raw := make(chan Product)
	ready := make(chan Product)
	var wg sync.WaitGroup
	wg.Go(func() { parse(a, raw) })
	wg.Go(func() { parse(b, raw) })
	go func() {
		wg.Wait()
		close(raw)
	}()
	stocked := make(chan Product)
	go inStock(raw, stocked)
	go unique(stocked, ready)

	var products []Product
	for p := range ready {
		products = append(products, p)
	}
	return products
}

func main() {
	for _, p := range importAll(chorsu, oloy) {
		fmt.Println(p.Name, p.Stock)
	}
}
