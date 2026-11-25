package main

import (
	"fmt"
	"sync"
)

// Guru: mutex'siz. Har bir goroutine faqat oʻz map'iga yozadi,
// wg.Wait() dan keyin bitta goroutine (main) ularni birlashtiradi.
func main() {
	qismlar := make([]map[int]int, 2)
	var wg sync.WaitGroup
	for g := range 2 {
		wg.Go(func() {
			m := map[int]int{}
			for i := range 1_000_000 {
				m[i] = g
			}
			qismlar[g] = m // har bir goroutine slice'ning oʻz elementiga yozadi
		})
	}
	wg.Wait()

	natija := map[int]int{}
	for _, m := range qismlar {
		for k, v := range m {
			natija[k] = v
		}
	}
	fmt.Println("tugadi:", len(natija))
}
