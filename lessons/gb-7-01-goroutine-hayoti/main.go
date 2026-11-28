// bozor import'ining birinchi qadami (GB-7.01): bir nechta yetkazib beruvchidan narx roʻyxatini
// bir vaqtda oʻqish. Tarmoq oʻrniga time.Sleep: har yetkazuvchi bir soniyada javob beradi.
// Toʻliq importer (cmd/import, bekor qilish, chegaralangan worker'lar) GB-7.07'da.
package main

import (
	"fmt"
	"time"
)

// offers — har yetkazib beruvchining narx roʻyxatida nechta mahsulot bor (laboratoriya maʼlumoti).
var offers = map[string]int{"chorsu": 40, "oloy": 120, "siyob": 5}

// fetch — bitta yetkazuvchining roʻyxatini "oʻqiydi" va undagi mahsulotlar sonini qaytaradi.
func fetch(name string) int {
	fmt.Println("boshlandi:", name)
	time.Sleep(time.Second)
	fmt.Println("tayyor:", name)
	return offers[name]
}

// fetchAll har yetkazuvchi uchun fetch'ni chaqiradi va natijalarni xuddi suppliers tartibida qaytaradi:
// counts[i] = fetch(suppliers[i]).
//
// VAZIFA (TASK.md, 1-qadam): hozir fetchAll ketma-ket ishlaydi (3 yetkazuvchi = 3 soniya). Uni videodagi
// main'dagidek qayta yozing: har fetch alohida goroutine'da (wg.Go), har goroutine oʻz counts[i]'iga yozadi,
// fetchAll esa qaytishdan oldin hammasini kutadi (wg.Wait). fetchAll — bu goroutine'larning egasi.
func fetchAll(suppliers []string, fetch func(string) int) []int {
	counts := make([]int, len(suppliers))
	for i, name := range suppliers {
		counts[i] = fetch(name)
	}
	return counts
}

// fetchAllAddDone — 🟡 VAZIFA, ixtiyoriy (TASK.md, 2-qadam): xuddi fetchAll, lekin eski shaklda yozilgan:
// wg.Add(1), go func() { defer wg.Done() … }(). Ishdagi koʻp kod shunday. Yozmaguningizcha nil
// qaytaradi va uning tekshiruvi SKIP boʻladi.
func fetchAllAddDone(suppliers []string, fetch func(string) int) []int {
	return nil
}

func main() {
	start := time.Now()
	counts := fetchAll([]string{"chorsu", "oloy", "siyob"}, fetch)
	fmt.Println("mahsulotlar:", counts)
	fmt.Println("jami:", time.Since(start))
}
