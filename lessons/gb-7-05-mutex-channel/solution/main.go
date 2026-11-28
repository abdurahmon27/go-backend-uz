package main

import (
	"fmt"
	"sync"
	"sync/atomic"
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

// sold — sanagich: sakkiz goroutine bir vaqtda oshiradi, qulfsiz.
var sold atomic.Int64

func main() {
	locked := &LockedStore{stock: map[int64]int{1: 800_000}}
	measure("mutex", func() { locked.Sell(1, 1) })
	owned := NewOwnedStore(map[int64]int{1: 800_000})
	measure("egasi", func() { owned.Sell(1, 1) })
	owned.Close()
	measure("atomic", func() { sold.Add(1) })
	fmt.Println("sotildi:", sold.Load())
}
