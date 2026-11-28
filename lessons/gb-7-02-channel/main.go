// bozor importeri: ikki yetkazib beruvchining narx roʻyxatini channel'lar orqali bosqichma-bosqich oʻtkazadi (GB-7.02).
//
//	parse (chorsu) ─┐
//	                ├─ raw ─ inStock ─ ready ─ importAll
//	parse (oloy)  ──┘
//
// Videoning oxiridagi kod shu; vazifa (TASK.md): unique bosqichi.
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
	// VAZIFA 1 (TASK.md): har ID'ni faqat bir marta yuboring,
	// oxirida out'ni yoping.
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
	go inStock(raw, ready)
	// VAZIFA 2 (TASK.md): inStock va natija orasiga unique bosqichini qoʻying.

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
