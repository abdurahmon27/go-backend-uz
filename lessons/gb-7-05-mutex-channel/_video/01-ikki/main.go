package main

import (
	"fmt"
	"sync"
	"time"
)

// measure op'ni sakkiz goroutine'dan jami 800 000 marta chaqiradi va bitta chaqiruvga ketgan oʻrtacha vaqtni chiqaradi.
func measure(name string, op func()) {
	start := time.Now()
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 100_000 {
				op()
			}
		})
	}
	wg.Wait()
	fmt.Println(name, time.Since(start)/800_000)
}

func main() {
	locked := &LockedStore{stock: map[int64]int{1: 800_000}}
	measure("mutex", func() { locked.Sell(1, 1) })
	owned := NewOwnedStore(map[int64]int{1: 800_000})
	measure("egasi", func() { owned.Sell(1, 1) })
	owned.Close()
}
