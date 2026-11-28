package main

import (
	"fmt"
	"sync"
	"time"
)

// measure — 01-ikki'dagidek, lekin goroutine'lar soni g: 800 000 chaqiruv g ta goroutine'ga boʻlinadi.
// Natija: jami vaqt / 800 000, ya'ni bitta chaqiruvga toʻgʻri keladigan vaqt (g = 1 da bu bitta chaqiruvning oʻzi).
func measure(name string, g int, op func()) {
	start := time.Now()
	var wg sync.WaitGroup
	for range g {
		wg.Go(func() {
			for range 800_000 / g {
				op()
			}
		})
	}
	wg.Wait()
	fmt.Println(name, g, time.Since(start)/800_000)
}

func main() {
	for _, g := range []int{1, 8} {
		locked := &LockedStore{stock: map[int64]int{1: 800_000}}
		measure("mutex", g, func() { locked.Sell(1, 1) })
		owned := NewOwnedStore(map[int64]int{1: 800_000})
		measure("egasi", g, func() { owned.Sell(1, 1) })
		owned.Close()
	}
}
