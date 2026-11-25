package main

import (
	"fmt"
	"sync"
)

func main() {
	m := map[int]int{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	for g := range 2 {
		wg.Go(func() {
			for i := range 1_000_000 {
				mu.Lock()
				m[i] = g
				mu.Unlock()
			}
		})
	}
	wg.Wait()
	fmt.Println("tugadi:", len(m))
}
