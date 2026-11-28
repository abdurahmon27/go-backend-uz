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

var stock atomic.Int64

// sellOne: CompareAndSwap — qiymat hali biz oʻqigan old boʻlsa, old-1 yozadi; boshqa goroutine
// oraliqda oʻzgartirgan boʻlsa, false qaytaradi va biz qaytadan oʻqiymiz.
func sellOne() bool {
	for {
		old := stock.Load()
		if old <= 0 {
			return false
		}
		if stock.CompareAndSwap(old, old-1) {
			return true
		}
	}
}

func main() {
	stock.Store(100_000)
	measure("cas", func() { sellOne() })
	fmt.Println("qoldiq:", stock.Load())
}
