// bozor import'ining birinchi qadami (GB-7.01): bir nechta yetkazib beruvchidan narx roʻyxatini
// bir vaqtda oʻqish. Tarmoq oʻrniga time.Sleep: har yetkazuvchi bir soniyada javob beradi.
// Toʻliq importer (cmd/import, bekor qilish, chegaralangan worker'lar) GB-7.07'da.
package main

import (
	"fmt"
	"sync"
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
// Javob (TASK.md, 1-qadam): har fetch alohida goroutine'da (wg.Go), har goroutine oʻz counts[i]'iga yozadi
// (har element alohida xotira joyi: data race yoʻq), fetchAll qaytishdan oldin hammasini kutadi (wg.Wait).
// fetchAll — bu goroutine'larning egasi: u qaytganda ulardan birortasi ham ishlamayapti.
func fetchAll(suppliers []string, fetch func(string) int) []int {
	counts := make([]int, len(suppliers))
	var wg sync.WaitGroup
	for i, name := range suppliers {
		wg.Go(func() { counts[i] = fetch(name) })
	}
	wg.Wait()
	return counts
}

// fetchAllAddDone — 🟡 javob (TASK.md, 2-qadam): xuddi fetchAll, eski shaklda. wg.Add(1) go'dan OLDIN
// (goroutine ichida emas: aks holda Wait hisob hali nol paytda qaytib ketishi mumkin), Done esa defer bilan.
// wg.Go shu ikki qatorni oʻzi bajaradi, bitta farq bilan: f panic qilsa, wg.Go Done'ni chaqirmaydi va panic'ni
// qayta tashlaydi, Wait ochilib, dastur panic tugashidan oldin chiqib ketmasin deb (src/sync/waitgroup.go).
// Hujjat yangi kodda wg.Go'ni tavsiya qiladi (go doc sync.WaitGroup.Add).
func fetchAllAddDone(suppliers []string, fetch func(string) int) []int {
	counts := make([]int, len(suppliers))
	var wg sync.WaitGroup
	for i, name := range suppliers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counts[i] = fetch(name)
		}()
	}
	wg.Wait()
	return counts
}

func main() {
	start := time.Now()
	counts := fetchAll([]string{"chorsu", "oloy", "siyob"}, fetch)
	fmt.Println("mahsulotlar:", counts)
	fmt.Println("jami:", time.Since(start))
}
