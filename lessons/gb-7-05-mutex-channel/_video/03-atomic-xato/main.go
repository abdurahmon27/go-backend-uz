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

// stock — oʻrik qoldigʻi, 100 000 dona. Xato: tekshirish va kamaytirish — ikki alohida atomik amal.
var stock atomic.Int64

func main() {
	stock.Store(100_000)
	measure("atomic", func() {
		if stock.Load() > 0 {
			stock.Add(-1)
		}
	})
	fmt.Println("qoldiq:", stock.Load())
}
