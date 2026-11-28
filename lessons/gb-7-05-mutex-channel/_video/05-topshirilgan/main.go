package main

import (
	"fmt"
	"sync"
)

// Xato: map loop'ga topshirilgandan keyin main uni oʻzi oʻqiydi. Endi map'ga ikki goroutine tegadi.
func main() {
	stock := map[int64]int{1: 1000}
	owned := NewOwnedStore(stock)
	var wg sync.WaitGroup
	for range 4 {
		wg.Go(func() {
			for range 100 {
				owned.Sell(1, 1)
			}
		})
	}
	fmt.Println("qoldiq:", stock[1]) // topshirilgan map'ga tegdik
	wg.Wait()
	owned.Close()
}
